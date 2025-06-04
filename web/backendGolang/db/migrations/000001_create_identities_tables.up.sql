CREATE TABLE appuser (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    no_telepon VARCHAR(20) NOT NULL,
    alamat TEXT,
    profile VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ(6),
    updated_at TIMESTAMPTZ(6)
);

CREATE TABLE data (
    id SERIAL PRIMARY KEY,
    jeniskendaraan VARCHAR(255) NOT NULL,
    jenispelanggaran VARCHAR(255) NOT NULL,
    lokasi VARCHAR(20) NOT NULL,
    waktu DATE,
    kecepatan VARCHAR(255) NOT NULL,
    gambar VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ(6),
    updated_at TIMESTAMPTZ(6)
);