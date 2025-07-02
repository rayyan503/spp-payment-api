# Backend API - Sistem Pembayaran SPP Online

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)
![Framework](https://img.shields.io/badge/Framework-Gin-green.svg)
![ORM](https://img.shields.io/badge/ORM-GORM-orange.svg)
![Database](https://img.shields.io/badge/Database-MySQL-blue.svg)
![License](https://img.shields.io/badge/License-MIT-lightgrey.svg)

Ini adalah backend API untuk aplikasi **Sistem Pembayaran SPP Online** tingkat Sekolah Dasar (SD). Proyek ini dibangun menggunakan Go (Golang) dengan framework Gin untuk performa tinggi dan GORM sebagai ORM untuk interaksi database yang efisien.

## Fitur Utama

-   **Manajemen Pengguna & Peran**: Sistem otentikasi berbasis peran (Admin, Bendahara, Siswa) menggunakan JWT.
-   **Akses Terkontrol**: Endpoint diamankan berdasarkan peran pengguna menggunakan middleware.
-   **Manajemen Data Master**: Pengelolaan data inti seperti tingkat kelas, biaya SPP, dan data kelas oleh Admin.
-   **Manajemen Siswa**: Fungsionalitas CRUD lengkap untuk data siswa oleh Bendahara.
-   **Siklus Penagihan**: Pembuatan periode tagihan SPP bulanan/tahunan.
-   **Generator Tagihan Otomatis**: Kemampuan untuk membuat tagihan SPP secara massal untuk semua siswa aktif dalam satu klik.
-   **Integrasi Payment Gateway**: Terhubung dengan **Midtrans** untuk memproses pembayaran online.
-   **Notifikasi Real-time**: Penanganan notifikasi (webhook) dari Midtrans untuk memperbarui status pembayaran secara otomatis.
-   **Pelaporan**: Laporan keuangan sederhana untuk memantau status pembayaran per kelas, per siswa, dan keseluruhan.

## Teknologi yang Digunakan

-   **Bahasa**: Go (v1.21+)
-   **Web Framework**: Gin Gonic
-   **ORM**: GORM
-   **Database**: MySQL
-   **Otentikasi**: JSON Web Tokens (JWT)
-   **Payment Gateway**: Midtrans
-   **Manajemen Dependensi**: Go Modules
-   **Konfigurasi**: Viper & Dotenv

## Prasyarat

-   [Go](https://golang.org/dl/) versi 1.21 atau lebih tinggi
-   [MySQL](https://www.mysql.com/downloads/)
-   [Git](https://git-scm.com/downloads/)

## Instalasi & Menjalankan Lokal

Ikuti langkah-langkah berikut untuk menjalankan proyek ini di lingkungan lokal Anda.

1.  **Clone Repository**
    ```sh
    git clone https://github.com/HIUNCY/spp-payment-api
    cd spp-payment-api
    ```

2.  **Setup Database**
    -   Buat sebuah database baru di MySQL Anda (misalnya, `spp-payment`).
    -   Impor skema dan data awal dari file `spp.sql` yang ada di repository.
        ```sh
        mysql -u [username] -p spp_sekolah < spp.sql
        ```

3.  **Konfigurasi Environment**
    -   Salin file `.env.example` menjadi `.env`.
        ```sh
        cp .env.example .env
        ```
    -   Buka file `.env` dan sesuaikan nilainya dengan konfigurasi lokal Anda, terutama untuk koneksi database dan kunci Midtrans.
        ```env
        # Server Configuration
        SERVER_PORT=8080

        # Database Configuration
        DB_HOST=localhost
        DB_PORT=3306
        DB_USER=root
        DB_PASSWORD=password
        DB_NAME=spp-payment

        # JWT Configuration
        JWT_SECRET_KEY=ini_rahasia_banget_jangan_disebar
        JWT_EXPIRATION_HOURS=72

        # Midtrans Configuration
        MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxxxxxxxxxxxxxx
        MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxxxxxxxxxxxxxx
        MIDTRANS_ENVIRONMENT=sandbox
        ```

4.  **Install Dependensi**
    ```sh
    go mod tidy
    ```

5.  **Jalankan Aplikasi**
    -   Server akan berjalan di port yang ditentukan di file `.env` (default: 8080).
    ```sh
    go run main.go
    ```

## Struktur Proyek

Proyek ini menggunakan arsitektur berlapis (*Layered Architecture*) untuk memisahkan tanggung jawab dan menjaga kode agar tetap bersih dan *maintainable*.
```
internal/
├── config/       # Manajemen konfigurasi (env)
├── handler/      # Layer presentasi (HTTP handlers, routing)
├── middleware/   # Middleware untuk otentikasi, logging, dll.
├── model/        # Struct data yang merepresentasikan tabel database (GORM models)
├── repository/   # Layer akses data (semua query database)
├── service/      # Layer logika bisnis
└── utils/        # Fungsi-fungsi bantuan (JWT, hashing, response)
```

## Dokumentasi API

Berikut adalah dokumentasi untuk endpoint yang telah diimplementasikan.

<details>
<summary><b>Otentikasi</b></summary>

### Login Pengguna
-   `POST /api/v1/login`
-   **Otorisasi**: Publik
-   **Request Body**:
    ```json
    {
        "email": "admin@sekolah.sch.id",
        "password": "password"
    }
    ```
-   **Response Sukses (200 OK)**:
    ```json
    {
        "status": "success",
        "message": "Login berhasil",
        "data": {
            "token": "jwt.token.string"
        }
    }
    ```

### Mendapatkan Profil Pengguna Login
-   `GET /api/v1/me`
-   **Otorisasi**: Admin, Bendahara, Siswa
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Response Sukses (200 OK)**:
    ```json
    {
        "status": "success",
        "message": "Profil pengguna berhasil diambil",
        "data": {
            "id": 1,
            "nama_lengkap": "Administrator",
            "email": "admin@sekolah.sch.id",
            "role": "admin"
        }
    }
    ```

</details>

<details>
<summary><b>Admin - Manajemen Pengguna</b></summary>

### Membuat Pengguna Baru
-   `POST /api/v1/admin/users`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "nama_lengkap": "Bendahara Sekolah",
        "email": "bendahara@sekolah.sch.id",
        "password": "password123",
        "role_id": 2
    }
    ```

### Mendapatkan Daftar Pengguna
-   `GET /api/v1/admin/users`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Query Params (Opsional)**:
    -   `page` (angka), `limit` (angka), `role_id` (angka), `search` (string)

### Mendapatkan Detail Pengguna
-   `GET /api/v1/admin/users/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Pengguna
-   `PUT /api/v1/admin/users/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "nama_lengkap": "Bendahara Utama Update",
        "email": "bendahara.utama@sekolah.sch.id",
        "role_id": 2
    }
    ```

### Menghapus Pengguna
-   `DELETE /api/v1/admin/users/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

<details>
<summary><b>Admin - Manajemen Tingkat Kelas</b></summary>

### Membuat Tingkat Kelas Baru
-   `POST /api/v1/admin/class-levels`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tingkat": 1,
        "nama_tingkat": "Kelas 1",
        "biaya_spp": 150000
    }
    ```

### Mendapatkan Semua Tingkat Kelas
-   `GET /api/v1/admin/class-levels`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Mendapatkan Detail Tingkat Kelas
-   `GET /api/v1/admin/class-levels/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Tingkat Kelas
-   `PUT /api/v1/admin/class-levels/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tingkat": 1,
        "nama_tingkat": "Kelas 1",
        "biaya_spp": 155000,
        "status": "aktif"
    }
    ```

### Menghapus Tingkat Kelas
-   `DELETE /api/v1/admin/class-levels/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

<details>
<summary><b>Admin - Manajemen Kelas</b></summary>

### Membuat Kelas Baru
-   `POST /api/v1/admin/classes`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tingkat_id": 1,
        "nama_kelas": "1A",
        "wali_kelas": "Bu Sari",
        "kapasitas": 30
    }
    ```

### Mendapatkan Semua Kelas
-   `GET /api/v1/admin/classes`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Mendapatkan Detail Kelas
-   `GET /api/v1/admin/classes/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Kelas
-   `PUT /api/v1/admin/classes/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tingkat_id": 1,
        "nama_kelas": "1A",
        "wali_kelas": "Sari Hartati, S.Pd.",
        "kapasitas": 32,
        "status": "aktif"
    }
    ```

### Menghapus Kelas
-   `DELETE /api/v1/admin/classes/{id}`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

<details>
<summary><b>Admin - Manajemen Pengaturan Sistem</b></summary>

### Mendapatkan Semua Pengaturan
-   `GET /api/v1/admin/settings`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Response Sukses (200 OK)**:
    ```json
    {
        "status": "success",
        "message": "Data pengaturan berhasil diambil",
        "data": {
            "nama_sekolah": "SD Negeri 1 Contoh",
            "alamat_sekolah": "Jl. Pendidikan No. 1, Kota",
            "telepon_sekolah": "021-1234567",
            "email_sekolah": "info@sekolah.sch.id",
            "tahun_ajaran_aktif": "2024/2025",
            "midtrans_server_key": "",
            "midtrans_client_key": "",
            "midtrans_environment": "sandbox"
        }
    }
    ```

### Memperbarui Pengaturan
-   `PUT /api/v1/admin/settings`
-   **Otorisasi**: Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "nama_sekolah": "SDN Koding Jaya",
        "telepon_sekolah": "021-111222"
    }
    ```
-   **Response Sukses (200 OK)**:
    ```json
    {
        "status": "success",
        "message": "Pengaturan berhasil diperbarui",
        "data": null
    }
    ```

</details>

<details>
<summary><b>Bendahara - Manajemen Siswa</b></summary>

### Membuat Siswa Baru
-   `POST /api/v1/treasurer/students`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "email": "budi.siswa@email.com",
        "password": "password123",
        "nisn": "1234567890",
        "kelas_id": 1,
        "nama_lengkap": "Budi Santoso",
        "jenis_kelamin": "L",
        "tempat_lahir": "Jakarta",
        "tanggal_lahir": "2018-07-15",
        "alamat": "Jl. Merdeka No. 10",
        "nama_orangtua": "Bapak Santoso",
        "telepon_orangtua": "08123456789",
        "tahun_masuk": 2024
    }
    ```

### Mendapatkan Daftar Siswa
-   `GET /api/v1/treasurer/students`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Query Params (Opsional)**:
    -   `page` (angka): Nomor halaman.
    -   `limit` (angka): Jumlah data per halaman.
    -   `kelas_id` (angka): Filter berdasarkan ID kelas.
    -   `search` (string): Cari berdasarkan nama atau NISN.

### Mendapatkan Detail Siswa
-   `GET /api/v1/treasurer/students/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Siswa
-   `PUT /api/v1/treasurer/students/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "nisn": "1234567890",
        "kelas_id": 2,
        "nama_lengkap": "Budi Santoso",
        "jenis_kelamin": "L",
        "tempat_lahir": "Jakarta",
        "tanggal_lahir": "2018-07-15",
        "alamat": "Jl. Kemerdekaan No. 12",
        "nama_orangtua": "Bapak Santoso",
        "telepon_orangtua": "08123456789",
        "tahun_masuk": 2024,
        "status": "aktif",
        "email": "budi.santoso@email.com",
        "status_user": "aktif"
    }
    ```

### Menghapus Siswa
-   `DELETE /api/v1/treasurer/students/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

<details>
<summary><b>Bendahara - Manajemen Periode SPP</b></summary>

### Membuat Periode Baru
-   `POST /api/v1/treasurer/periods`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tahun_ajaran": "2025/2026",
        "bulan": 7,
        "nama_bulan": "Juli",
        "tanggal_mulai": "2025-07-01",
        "tanggal_selesai": "2025-07-31"
    }
    ```

### Mendapatkan Daftar Periode
-   `GET /api/v1/treasurer/periods`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Query Params (Opsional)**:
    -   `tahun_ajaran` (string): Filter berdasarkan tahun ajaran, contoh: `2025/2026`.

### Mendapatkan Detail Periode
-   `GET /api/v1/treasurer/periods/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Periode
-   `PUT /api/v1/treasurer/periods/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "tahun_ajaran": "2025/2026",
        "bulan": 7,
        "nama_bulan": "Juli",
        "tanggal_mulai": "2025-07-01",
        "tanggal_selesai": "2025-08-10",
        "status": "aktif"
    }
    ```

### Menghapus Periode
-   `DELETE /api/v1/treasurer/periods/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

<details>
<summary><b>Bendahara - Manajemen Tagihan</b></summary>

### Generate Tagihan Massal
-   `POST /api/v1/treasurer/periods/{id}/generate-bills`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Fungsi**: Membuat tagihan SPP untuk semua siswa aktif berdasarkan ID periode yang diberikan.

### Mendapatkan Daftar Tagihan
-   `GET /api/v1/treasurer/bills`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Query Params (Opsional)**:
    -   `page` (angka): Nomor halaman.
    -   `limit` (angka): Jumlah data per halaman.
    -   `periode_id` (angka): Filter berdasarkan ID periode.
    -   `siswa_id` (angka): Filter berdasarkan ID siswa.
    -   `status_pembayaran` (string): Filter berdasarkan status (`belum_bayar`, `pending`, `lunas`).

### Mendapatkan Detail Tagihan
-   `GET /api/v1/treasurer/bills/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

### Memperbarui Tagihan (Manual)
-   `PUT /api/v1/treasurer/bills/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`
-   **Request Body**:
    ```json
    {
        "jumlah_tagihan": 160000,
        "status_pembayaran": "lunas"
    }
    ```

### Menghapus Tagihan (Manual)
-   `DELETE /api/v1/treasurer/bills/{id}`
-   **Otorisasi**: Bendahara, Admin
-   **Header**: `Authorization: Bearer <TOKEN>`

</details>

## Kontribusi

Kontribusi dalam bentuk *pull request*, isu, atau ide fitur sangat diterima.

## Lisensi

Proyek ini dilisensikan di bawah [Lisensi MIT](LICENSE).
