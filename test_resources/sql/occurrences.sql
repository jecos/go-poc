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