-- Database SPP Online SD
-- Sistem Pembayaran SPP dengan Midtrans Payment Gateway

-- Tabel untuk menyimpan data tingkat kelas dan biaya SPP
CREATE TABLE tingkat_kelas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tingkat INT NOT NULL COMMENT 'Tingkat kelas (1-6)',
    nama_tingkat VARCHAR(50) NOT NULL COMMENT 'Nama tingkat (Kelas 1, Kelas 2, dst)',
    biaya_spp DECIMAL(12,2) NOT NULL COMMENT 'Biaya SPP per bulan',
    status ENUM('aktif', 'nonaktif') DEFAULT 'aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tingkat (tingkat)
);

-- Tabel untuk menyimpan data kelas (1A, 1B, 2A, 2B, dst)
CREATE TABLE kelas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tingkat_id INT NOT NULL,
    nama_kelas VARCHAR(10) NOT NULL COMMENT 'Nama kelas (1A, 1B, 2A, dst)',
    wali_kelas VARCHAR(100) DEFAULT NULL,
    kapasitas INT DEFAULT 30,
    status ENUM('aktif', 'nonaktif') DEFAULT 'aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tingkat_id) REFERENCES tingkat_kelas(id) ON DELETE CASCADE,
    UNIQUE KEY unique_nama_kelas (nama_kelas)
);

-- Tabel untuk role/peran pengguna
CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nama_role VARCHAR(50) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_role (nama_role)
);

-- Tabel untuk menyimpan data pengguna (admin, bendahara, siswa)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
    nama_lengkap VARCHAR(100) NOT NULL,
    status ENUM('aktif', 'nonaktif') DEFAULT 'aktif',
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,
    INDEX idx_email (email),
    INDEX idx_role (role_id)
);

-- Tabel untuk menyimpan data siswa
CREATE TABLE siswa (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    nisn VARCHAR(20) NOT NULL UNIQUE,
    kelas_id INT NOT NULL,
    nama_lengkap VARCHAR(100) NOT NULL,
    jenis_kelamin ENUM('L', 'P') NOT NULL,
    tempat_lahir VARCHAR(50),
    tanggal_lahir DATE,
    alamat TEXT,
    nama_orangtua VARCHAR(100),
    telepon_orangtua VARCHAR(20),
    tahun_masuk YEAR,
    status ENUM('aktif', 'pindah', 'lulus', 'keluar') DEFAULT 'aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE RESTRICT,
    INDEX idx_nisn (nisn),
    INDEX idx_kelas (kelas_id),
    INDEX idx_status (status)
);

-- Tabel untuk mengatur periode pembayaran SPP
CREATE TABLE periode_spp (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tahun_ajaran VARCHAR(20) NOT NULL COMMENT 'Format: 2024/2025',
    bulan INT NOT NULL COMMENT 'Bulan SPP (1-12)',
    nama_bulan VARCHAR(20) NOT NULL COMMENT 'Nama bulan (Januari, Februari, dst)',
    tanggal_mulai DATE NOT NULL COMMENT 'Tanggal mulai pembayaran',
    tanggal_selesai DATE NOT NULL COMMENT 'Batas akhir pembayaran',
    status ENUM('belum_aktif', 'aktif', 'selesai') DEFAULT 'belum_aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_periode (tahun_ajaran, bulan),
    INDEX idx_status (status),
    INDEX idx_tanggal (tanggal_mulai, tanggal_selesai)
);

-- Tabel untuk menyimpan tagihan SPP per siswa per periode
CREATE TABLE tagihan_spp (
    id INT PRIMARY KEY AUTO_INCREMENT,
    siswa_id INT NOT NULL,
    periode_id INT NOT NULL,
    jumlah_tagihan DECIMAL(12,2) NOT NULL,
    status_pembayaran ENUM('belum_bayar', 'pending', 'lunas') DEFAULT 'belum_bayar',
    tanggal_jatuh_tempo DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE,
    FOREIGN KEY (periode_id) REFERENCES periode_spp(id) ON DELETE CASCADE,
    UNIQUE KEY unique_tagihan (siswa_id, periode_id),
    INDEX idx_status (status_pembayaran),
    INDEX idx_jatuh_tempo (tanggal_jatuh_tempo)
);

-- Tabel untuk menyimpan data pembayaran dan integrasi dengan Midtrans
CREATE TABLE pembayaran (
    id INT PRIMARY KEY AUTO_INCREMENT,
    tagihan_id INT NOT NULL,
    siswa_id INT NOT NULL,
    order_id VARCHAR(100) NOT NULL UNIQUE COMMENT 'Order ID untuk Midtrans',
    transaction_id VARCHAR(100) NULL COMMENT 'Transaction ID dari Midtrans',
    jumlah_bayar DECIMAL(12,2) NOT NULL,
    metode_pembayaran VARCHAR(50) NULL COMMENT 'bank_transfer, e_wallet, credit_card, dll',
    status_pembayaran ENUM('pending', 'settlement', 'cancel', 'expire', 'failure') DEFAULT 'pending',
    tanggal_pembayaran TIMESTAMP NULL,
    tanggal_settlement TIMESTAMP NULL,
    midtrans_response JSON NULL COMMENT 'Response lengkap dari Midtrans',
    keterangan TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tagihan_id) REFERENCES tagihan_spp(id) ON DELETE CASCADE,
    FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_status (status_pembayaran),
    INDEX idx_tanggal_bayar (tanggal_pembayaran)
);

-- Tabel untuk log aktivitas sistem
CREATE TABLE log_aktivitas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NULL,
    aktivitas VARCHAR(255) NOT NULL,
    detail TEXT NULL,
    ip_address VARCHAR(45) NULL,
    user_agent TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user (user_id),
    INDEX idx_created_at (created_at)
);

-- Tabel untuk pengaturan sistem
CREATE TABLE pengaturan (
    id INT PRIMARY KEY AUTO_INCREMENT,
    key_setting VARCHAR(100) NOT NULL UNIQUE,
    value_setting TEXT NOT NULL,
    deskripsi TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ============================
-- INSERT DATA AWAL
-- ============================

-- Insert roles
INSERT INTO roles (nama_role, deskripsi) VALUES
('admin', 'Administrator sistem'),
('bendahara', 'Bendahara sekolah'),
('siswa', 'Siswa sekolah');

-- Insert tingkat kelas dengan biaya SPP (sesuaikan dengan kebutuhan)
INSERT INTO tingkat_kelas (tingkat, nama_tingkat, biaya_spp) VALUES
(1, 'Kelas 1', 150000.00),
(2, 'Kelas 2', 155000.00),
(3, 'Kelas 3', 160000.00),
(4, 'Kelas 4', 165000.00),
(5, 'Kelas 5', 170000.00),
(6, 'Kelas 6', 175000.00);

-- Insert kelas (2 kelas per tingkat: A dan B)
INSERT INTO kelas (tingkat_id, nama_kelas, wali_kelas) VALUES
(1, '1A', 'Bu Sari'),
(1, '1B', 'Bu Dewi'),
(2, '2A', 'Bu Rina'),
(2, '2B', 'Bu Maya'),
(3, '3A', 'Bu Fitri'),
(3, '3B', 'Bu Indri'),
(4, '4A', 'Pak Budi'),
(4, '4B', 'Pak Joko'),
(5, '5A', 'Bu Lestari'),
(5, '5B', 'Bu Wati'),
(6, '6A', 'Pak Ahmad'),
(6, '6B', 'Pak Hendra');

-- Insert user admin dan bendahara default
INSERT INTO users (email, password, role_id, nama_lengkap) VALUES
('admin@sekolah.sch.id', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 1, 'Administrator'),
('bendahara@sekolah.sch.id', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 2, 'Bendahara Sekolah');
-- Password default: "password" (hash bcrypt)

-- Insert pengaturan awal
INSERT INTO pengaturan (key_setting, value_setting, deskripsi) VALUES
('nama_sekolah', 'SD Negeri 1 Contoh', 'Nama sekolah'),
('alamat_sekolah', 'Jl. Pendidikan No. 1, Kota', 'Alamat sekolah'),
('telepon_sekolah', '021-1234567', 'Nomor telepon sekolah'),
('email_sekolah', 'info@sekolah.sch.id', 'Email resmi sekolah'),
('tahun_ajaran_aktif', '2024/2025', 'Tahun ajaran yang sedang aktif'),
('midtrans_server_key', '', 'Server Key Midtrans'),
('midtrans_client_key', '', 'Client Key Midtrans'),
('midtrans_environment', 'sandbox', 'Environment Midtrans (sandbox/production)');

-- ============================
-- VIEW UNTUK LAPORAN
-- ============================

-- View untuk laporan pembayaran per siswa
CREATE VIEW v_laporan_siswa AS
SELECT
    s.nisn,
    s.nama_lengkap,
    k.nama_kelas,
    tk.nama_tingkat,
    ps.tahun_ajaran,
    ps.nama_bulan,
    ts.jumlah_tagihan,
    ts.status_pembayaran,
    ts.tanggal_jatuh_tempo,
    p.tanggal_settlement,
    p.metode_pembayaran
FROM siswa s
JOIN kelas k ON s.kelas_id = k.id
JOIN tingkat_kelas tk ON k.tingkat_id = tk.id
JOIN tagihan_spp ts ON s.id = ts.siswa_id
JOIN periode_spp ps ON ts.periode_id = ps.id
LEFT JOIN pembayaran p ON ts.id = p.tagihan_id AND p.status_pembayaran = 'settlement'
WHERE s.status = 'aktif';

-- View untuk laporan pembayaran per kelas
CREATE VIEW v_laporan_kelas AS
SELECT
    k.nama_kelas,
    tk.nama_tingkat,
    ps.tahun_ajaran,
    ps.nama_bulan,
    COUNT(ts.id) as total_siswa,
    SUM(CASE WHEN ts.status_pembayaran = 'lunas' THEN 1 ELSE 0 END) as siswa_lunas,
    SUM(CASE WHEN ts.status_pembayaran = 'belum_bayar' THEN 1 ELSE 0 END) as siswa_belum_bayar,
    SUM(CASE WHEN ts.status_pembayaran = 'pending' THEN 1 ELSE 0 END) as siswa_pending,
    SUM(ts.jumlah_tagihan) as total_tagihan,
    SUM(CASE WHEN ts.status_pembayaran = 'lunas' THEN ts.jumlah_tagihan ELSE 0 END) as total_terbayar
FROM kelas k
JOIN tingkat_kelas tk ON k.tingkat_id = tk.id
JOIN siswa s ON k.id = s.kelas_id AND s.status = 'aktif'
JOIN tagihan_spp ts ON s.id = ts.siswa_id
JOIN periode_spp ps ON ts.periode_id = ps.id
GROUP BY k.id, ps.id;

-- View untuk laporan keseluruhan
CREATE VIEW v_laporan_keseluruhan AS
SELECT
    ps.tahun_ajaran,
    ps.nama_bulan,
    COUNT(ts.id) as total_tagihan,
    SUM(CASE WHEN ts.status_pembayaran = 'lunas' THEN 1 ELSE 0 END) as total_lunas,
    SUM(CASE WHEN ts.status_pembayaran = 'belum_bayar' THEN 1 ELSE 0 END) as total_belum_bayar,
    SUM(CASE WHEN ts.status_pembayaran = 'pending' THEN 1 ELSE 0 END) as total_pending,
    SUM(ts.jumlah_tagihan) as total_nominal_tagihan,
    SUM(CASE WHEN ts.status_pembayaran = 'lunas' THEN ts.jumlah_tagihan ELSE 0 END) as total_nominal_terbayar,
    ROUND((SUM(CASE WHEN ts.status_pembayaran = 'lunas' THEN ts.jumlah_tagihan ELSE 0 END) / SUM(ts.jumlah_tagihan)) * 100, 2) as persentase_pembayaran
FROM periode_spp ps
JOIN tagihan_spp ts ON ps.id = ts.periode_id
JOIN siswa s ON ts.siswa_id = s.id AND s.status = 'aktif'
GROUP BY ps.id;

-- ============================
-- STORED PROCEDURES
-- ============================

DELIMITER //

-- Procedure untuk generate tagihan SPP otomatis berdasarkan periode
CREATE PROCEDURE GenerateTagihanSPP(
    IN p_periode_id INT
)
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_siswa_id INT;
    DECLARE v_biaya_spp DECIMAL(12,2);
    DECLARE v_tanggal_jatuh_tempo DATE;

    -- Cursor untuk mengambil semua siswa aktif dengan biaya SPP
    DECLARE siswa_cursor CURSOR FOR
        SELECT s.id, tk.biaya_spp
        FROM siswa s
        JOIN kelas k ON s.kelas_id = k.id
        JOIN tingkat_kelas tk ON k.tingkat_id = tk.id
        WHERE s.status = 'aktif';

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    -- Ambil tanggal jatuh tempo dari periode
    SELECT tanggal_selesai INTO v_tanggal_jatuh_tempo
    FROM periode_spp
    WHERE id = p_periode_id;

    OPEN siswa_cursor;

    read_loop: LOOP
        FETCH siswa_cursor INTO v_siswa_id, v_biaya_spp;
        IF done THEN
            LEAVE read_loop;
        END IF;

        -- Insert tagihan jika belum ada
        INSERT IGNORE INTO tagihan_spp (siswa_id, periode_id, jumlah_tagihan, tanggal_jatuh_tempo)
        VALUES (v_siswa_id, p_periode_id, v_biaya_spp, v_tanggal_jatuh_tempo);

    END LOOP;

    CLOSE siswa_cursor;
END //

-- Procedure untuk update status pembayaran dari callback Midtrans
CREATE PROCEDURE UpdateStatusPembayaran(
    IN p_order_id VARCHAR(100),
    IN p_transaction_status VARCHAR(50),
    IN p_transaction_id VARCHAR(100),
    IN p_payment_type VARCHAR(50),
    IN p_settlement_time TIMESTAMP,
    IN p_midtrans_response JSON
)
BEGIN
    DECLARE v_tagihan_id INT;

    -- Update pembayaran
    UPDATE pembayaran
    SET
        transaction_id = p_transaction_id,
        status_pembayaran = p_transaction_status,
        metode_pembayaran = p_payment_type,
        tanggal_settlement = p_settlement_time,
        midtrans_response = p_midtrans_response,
        updated_at = CURRENT_TIMESTAMP
    WHERE order_id = p_order_id;

    -- Ambil tagihan_id
    SELECT tagihan_id INTO v_tagihan_id
    FROM pembayaran
    WHERE order_id = p_order_id;

    -- Update status tagihan
    IF p_transaction_status = 'settlement' THEN
        UPDATE tagihan_spp
        SET status_pembayaran = 'lunas'
        WHERE id = v_tagihan_id;
    ELSEIF p_transaction_status IN ('pending', 'challenge') THEN
        UPDATE tagihan_spp
        SET status_pembayaran = 'pending'
        WHERE id = v_tagihan_id;
    ELSE
        UPDATE tagihan_spp
        SET status_pembayaran = 'belum_bayar'
        WHERE id = v_tagihan_id;
    END IF;

END //

DELIMITER ;

-- ============================
-- INDEXES TAMBAHAN UNTUK PERFORMA
-- ============================

-- Index untuk pencarian cepat
CREATE INDEX idx_siswa_nama ON siswa(nama_lengkap);
CREATE INDEX idx_tagihan_status_periode ON tagihan_spp(status_pembayaran, periode_id);
CREATE INDEX idx_pembayaran_tanggal ON pembayaran(tanggal_settlement);

-- ============================
-- CONTOH PENGGUNAAN
-- ============================

/*
-- Contoh insert periode SPP untuk bulan Januari 2025
INSERT INTO periode_spp (tahun_ajaran, bulan, nama_bulan, tanggal_mulai, tanggal_selesai, status)
VALUES ('2024/2025', 1, 'Januari', '2025-01-01', '2025-01-31', 'aktif');

-- Generate tagihan untuk periode tersebut
CALL GenerateTagihanSPP(1);

-- Contoh query laporan per kelas
SELECT * FROM v_laporan_kelas WHERE tahun_ajaran = '2024/2025' AND nama_bulan = 'Januari';

-- Contoh query laporan keseluruhan
SELECT * FROM v_laporan_keseluruhan WHERE tahun_ajaran = '2024/2025';
*/
