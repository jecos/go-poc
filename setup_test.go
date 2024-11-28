package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	starrocksContainer testcontainers.Container
	containerStarted   bool
	once               sync.Once
)

func ParallelTestWithDb(t *testing.T, dbName string, testFunc func(t *testing.T, db *sql.DB)) {
	t.Parallel()
	db, dbName, err := initDb(dbName)
	if err != nil {
		fmt.Errorf("Failed to init db connection: %v\n", err)
		os.Exit(1)
	}
	testFunc(t, db)
	//Drop database
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
	if err != nil {
		fmt.Errorf("Failed to drop database: %v\n", err)
	}
	err = db.Close()
	if err != nil {
		fmt.Errorf("Failed to close db connection: %v\n", err)
		os.Exit(1)
	}
}
func initDb(folderName string) (*sql.DB, string, error) {
	ctx := context.Background()
	host, err := starrocksContainer.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container host: %v", err)
	}

	port, err := starrocksContainer.MappedPort(ctx, "9030")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("root:@tcp(%s:%s)/?interpolateParams=true", host, port.Port())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to StarRocks: %v", err)
	}

	// Create a database with the name of the folder and a random number
	dbName := fmt.Sprintf("%s_%d", folderName, rand.Intn(1000000))
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf("USE %s;", dbName))
	if err != nil {
		return nil, "", fmt.Errorf("failed to use database: %v", err)
	}

	// Read the list of .tsv files in the folder
	files, err := os.ReadDir(filepath.Join("test_resources", folderName))
	if err != nil {
		return nil, "", fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tsv" {
			err = createTableAndPopulateData(db, folderName, file)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create table and populate data: %v", err)
			}
		}
	}

	return db, dbName, nil
}

func createTableAndPopulateData(db *sql.DB, folderName string, file os.DirEntry) error {
	if filepath.Ext(file.Name()) != ".tsv" {
		return nil
	}

	tableName := strings.TrimSuffix(file.Name(), ".tsv")
	sqlFilePath := filepath.Join("test_resources", "sql", tableName+".sql")

	// Read the SQL file
	sqlFile, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Execute the SQL file to create the table
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Populate the table with data from the .tsv file
	tsvFilePath := filepath.Join("test_resources", folderName, file.Name())
	tsvFile, err := os.Open(tsvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open TSV file: %v", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)
	// Read the header line
	var columns []string
	if scanner.Scan() {
		header := scanner.Text()
		columns = strings.Split(header, "\t")
	}

	columnNames := strings.Join(columns, ", ")
	p := strings.Repeat("?, ", len(columns))
	p = p[:len(p)-2]
	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnNames, p)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, "\t")

		// Convert to interface{} for stmt.Exec
		args := make([]interface{}, len(values))
		for i, v := range values {
			args[i] = v // Database driver will handle type conversion
		}

		// Execute the prepared statement
		_, err = db.Exec(insertQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to insert row: %v", err)
		}
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error reading TSV file: %v", err)
	}

	return nil
}

func startStarRocksContainer() (testcontainers.Container, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "starrocks/allin1-ubuntu",
		ExposedPorts: []string{"9030/tcp", "8030/tcp", "8040/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("9030/tcp"),
			wait.ForListeningPort("8030/tcp"),
			wait.ForListeningPort("8040/tcp"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	defer time.Sleep(30 * time.Second)

	return container, nil
}

func setupContainer() {
	once.Do(func() {
		// Run the script to start the container if it's not already running
		cmd := exec.Command("docker", "ps", "-q", "-f", "name=starrocks_test")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Failed to check if StarRocks container is running: %v\n", err)
			os.Exit(1)
		}

		if len(output) == 0 {
			// Container is not running, start a new one
			containerStarted = true
			starrocksContainer, err = startStarRocksContainer()
			if err != nil {
				fmt.Printf("Failed to start StarRocks container: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Container is already running, attach to it
			fmt.Println("StarRocks container is already running.")
			ctx := context.Background()
			starrocksContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
				ContainerRequest: testcontainers.ContainerRequest{
					Name: "starrocks_test",
				},
				Started: true,
				Reuse:   true,
			})
			if err != nil {
				fmt.Printf("Failed to attach to StarRocks container: %v\n", err)
				os.Exit(1)
			}
		}
	})
}

func TestMain(m *testing.M) {

	setupContainer()

	code := m.Run()

	// Clean up

	if containerStarted {
		ctx := context.Background()
		starrocksContainer.Terminate(ctx)
	}

	os.Exit(code)
}
