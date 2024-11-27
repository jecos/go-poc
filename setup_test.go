package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
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
	db := initDb(dbName)
	testFunc(t, db)
	err := db.Close()
	if err != nil {
		return
	}
}
func initDb(folderName string) *sql.DB {
	ctx := context.Background()
	host, err := starrocksContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		os.Exit(1)
	}

	port, err := starrocksContainer.MappedPort(ctx, "9030")
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		os.Exit(1)
	}
	dsn := fmt.Sprintf("root:@tcp(%s:%s)/?interpolateParams=true", host, port.Port())
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		fmt.Printf("Failed to connect to StarRocks: %v\n", err)
		os.Exit(1)
	}

	// Create a database with the name of the folder and a random number
	dbName := fmt.Sprintf("%s_%d", folderName, rand.Intn(1000000))
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		os.Exit(1)
	}
	_, err = db.Exec(fmt.Sprintf("USE %s;", dbName))
	if err != nil {
		fmt.Printf("Failed to use database: %v\n", err)
		os.Exit(1)
	}

	// Create the occurrences table
	createTableSQL := `
  CREATE TABLE IF NOT EXISTS occurrences (
   seq_id INT,
   locus_id VARCHAR(255),
   quality BIGINT,
   filter VARCHAR(255),
   zygosity VARCHAR(255),
   pf DOUBLE,
   af DOUBLE,
   gnomad_v3_af DOUBLE,
   hgvsg VARCHAR(255),
   omim_inheritance_code VARCHAR(255),
   ad_ratio DOUBLE,
   variant_class VARCHAR(255),
   vep_impact VARCHAR(255),
   symbol VARCHAR(255),
   clinvar_interpretation VARCHAR(255),
   mane_select BOOLEAN,
   canonical BOOLEAN
  );
 `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("Failed to create table: %v\n", err)
		os.Exit(1)
	}

	// Populate the occurrences table with data from occurrence.tsv
	filePath := filepath.Join("test_resources", folderName, "occurrences.tsv")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Read the header line
	var columns []string
	if scanner.Scan() {
		header := scanner.Text()
		columns = strings.Split(header, "\t")
	}

	tableName := "occurrences"
	columnNames := strings.Join(columns, ", ")
	placeholders := strings.Repeat("?, ", len(columns))
	placeholders = placeholders[:len(placeholders)-2]
	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnNames, placeholders)
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
			log.Printf("Failed to insert row: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	return db
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
