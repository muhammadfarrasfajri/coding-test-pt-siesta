# Core Banking Transfer Engine - SIESTA Coding Test

Sebuah _RESTful API Service_ berbasis Golang untuk memproses transaksi transfer antar-rekening (_Core Banking_). Proyek ini dirancang dengan fokus utama pada **Data Integrity, Concurrency (Low Footprint), Atomicity,** dan **Clean Architecture**.

## 🚀 Fitur Unggulan (Core Engineering Decisions)

- **Deadlock Prevention (Lock Ordering):** Mencegah _deadlock_ saat terjadi transaksi silang (A -> B dan B -> A secara bersamaan) dengan mengimplementasikan `ORDER BY id FOR UPDATE` pada tingkat _database_. Row di-lock berdasarkan urutan leksikografis ID.
- **Atomicity & Double-Entry Bookkeeping:** Menggunakan transaksi manual (`sql.Tx`) di _Service Layer_ untuk memastikan proses debit, kredit, dan pencatatan _ledger_ mutasi berjalan secara atomik (_All-or-Nothing_).
- **Database-Level Idempotency:** Menjamin sistem kebal terhadap _Double Spending_ jika terjadi pengulangan _request_ (_retry_) dari klien akibat _timeout_. Pengecekan idempotensi mengandalkan _Unique Constraint_ PostgreSQL (`23505`) sehingga sangat _low footprint_ dan bebas _race condition_.
- **Clean Architecture & Modularity:** Pemisahan tugas yang tegas antara `Controller` (HTTP/Input validation), `Service` (Business Logic & Transaction Boundary), dan `Repository` (Raw SQL Query).
- **Custom Error Observability:** Menggunakan _custom error struct_ yang menyertakan `TraceID`, `ErrorCode`, dan `Severity` untuk mempermudah _debugging_ dan _tracing_ di _production_.
- **Table-Driven Unit Testing:** Pengujian logika _Service_ menggunakan _Mocking_ (`go-sqlmock` & `testify`) untuk mencapai _test coverage_ yang tinggi tanpa ketergantungan pada _database_ asli.

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
- **Migration & Seeding:** Goose (`github.com/pressly/goose/v3`)
- **Testing:** Testify & Go-SQLMock

## 📁 Struktur Direktori

```
.
├── bootstrap/               # Inisialisasi dependensi aplikasi
├── cmd/migrate/             # Entry point untuk script migrasi database (Goose)
├── controller/              # Layer presentasi (HTTP Handlers & JSON Binding)
├── db/                      # Konfigurasi dan koneksi database
├── migrations/              # Kumpulan file SQL untuk skema DB dan Seeder Dummy Data
├── model/                   # Domain entities dan format response (DTO)
├── pkg/                     # Utility dan Custom Error Definitions
├── repository/              # Layer Data Access (Raw SQL execution)
├── router/                  # Konfigurasi routing (Gin)
├── service/                 # Layer Business Logic (Transaksional DB)
├── .env.example             # Template environment variables
├── api_test.http            # Kumpulan HTTP Request untuk pengujian (REST Client)
└── main.go                  # Bootstrap utama aplikasi

```

## ⚙️ Persyaratan Sistem

- Go 1.20+
- PostgreSQL 13+

## 🚦 Cara Instalasi & Menjalankan

### 1. Konfigurasi Database

Buat database baru di PostgreSQL Anda dengan nama `core_banking_db`:

```
psql -U postgres -c "CREATE DATABASE core_banking_db;"

```

### 2. Setup Environment Variables

Salin file `.env.example` menjadi `.env`:

```
cp .env.example .env

```

Lalu buka file `.env` dan sesuaikan kredensial PostgreSQL Anda:

```
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASS=password_postgres_anda
DB_NAME=core_banking_db
DB_USE=postgres
PORT=8080

```

### 3. Jalankan Migrasi & Seeding Data Dummy

Aplikasi ini sudah dilengkapi dengan _Database Seeder_. Perintah di bawah ini akan otomatis membuat tabel-tabel yang dibutuhkan sekaligus menyuntikkan 2 akun dummy (`A001` dan `B002`):

```
go run ./cmd/migrate up

```

### 4. Jalankan Server

```
go run main.go

```

_Server akan berjalan di `http://localhost:8080_`

## 🧪 Pengujian (Testing)

### Unit Testing

Untuk menjalankan _unit test_ pada logika bisnis (_Service Layer_) dan memastikan kode tahan terhadap _race condition_:

```
go test ./... -v -race

```

or

```
go test ./... -v

```

_(Catatan: Flag `-race` memerlukan CGO/C Compiler 64-bit. Jika Anda menjalankan ini di environment tanpa C Compiler aktif, cukup gunakan `go test ./... -v`)_

### API Testing (REST Client)

Untuk kemudahan pengujian tanpa perlu melakukan _setup_ Postman, saya telah menyediakan file **`api_test.http`** di _root_ direktori.

Jika Anda menggunakan **VS Code** (dengan ekstensi _REST Client_) atau **JetBrains GoLand**, Anda cukup membuka file tersebut dan menekan tombol **"Send Request"** untuk menguji berbagai skenario (Transfer Normal, Idempotency, Saldo Tidak Cukup, Validasi Error, dll).

---

_Dibuat oleh Muhammad Farras Fajri untuk Coding Test PT SIESTA._

```

```
