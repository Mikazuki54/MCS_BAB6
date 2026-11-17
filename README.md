# Dokumentasi Lengkap Proyek REST API Card Bridge (`bab6`)

Proyek ini adalah aplikasi **RESTful API** berbasis bahasa pemrograman **Go (Golang)** menggunakan framework **Gin Gonic** dan database **PostgreSQL**. Proyek ini menerapkan pola arsitektur berlapis (**Layered Architecture** / Repository Pattern) untuk mengelola data kartu (`card`) yang mencakup operasi pembacaan (`GET`), penambahan (`POST`), dan penghapusan (`DELETE`), serta sistem migrasi database otomatis menggunakan **`sql-migrate`** dengan fitur **`embed`**.

---

## Daftar Isi

1. [Struktur Direktori Proyek](#1-struktur-direktori-proyek)
2. [Penjelasan Rinci File, Syntax, Nama Variabel, dan Fungsi](#2-penjelasan-rinci-file-syntax-nama-variabel-dan-fungsi)
   - [2.1. `go.mod` & `go.sum` (Konfigurasi Modul & Dependensi)](#21-gomod--gosum-konfigurasi-modul--dependensi)
   - [2.2. `entities/card_bridge_entity.go` (Model Data)](#22-entitiescard_bridge_entitygo-model-data)
   - [2.3. `databases/sql_migrations/1_initiate.sql` (Skrip Migrasi SQL)](#23-databasessql_migrations1_initiatesql-skrip-migrasi-sql)
   - [2.4. `databases/database.go` (Koneksi & Migrasi Otomatis)](#24-databasesdatabasego-koneksi--migrasi-otomatis)
   - [2.5. `repositories/card_bridge_repo.go` (Akses Database / Query SQL)](#25-repositoriescard_bridge_repogo-akses-database--query-sql)
   - [2.6. `controllers/card_bridge_controllers.go` (Logika Bisnis & HTTP Handler)](#26-controllerscard_bridge_controllersgo-logika-bisnis--http-handler)
   - [2.7. `routers/card_bridge_routers.go` (Routing & CORS Middleware)](#27-routerscard_bridge_routersgo-routing--cors-middleware)
   - [2.8. `main.go` (Entry Point / Titik Masuk Utama)](#28-maingo-entry-point--titik-masuk-utama)
3. [Dokumentasi Endpoint API](#3-dokumentasi-endpoint-api)
4. [Cara Menjalankan Proyek](#4-cara-menjalankan-proyek)

---

## 1. Struktur Direktori Proyek

```text
testbab6/
├── controllers/
│   └── card_bridge_controllers.go      # Menangani request HTTP, parameter URL, dan response JSON
├── databases/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql              # Skrip SQL untuk membuat tabel 'card'
│   └── database.go                     # Manajemen koneksi DB & eksekusi migrasi embedded sql
├── entities/
│   └── card_bridge_entity.go           # Definisi struktur data struct 'Card'
├── repositories/
│   └── card_bridge_repo.go             # Eksekusi query SQL (SELECT, INSERT, DELETE) ke database
├── routers/
│   └── card_bridge_routers.go          # Definisi routing endpoint API dan konfigurasi CORS
├── go.mod                              # Identitas modul Go dan daftar package eksternal
├── go.sum                              # Checksum kriptografi dependensi
├── main.go                             # Entry point utama aplikasi Go
└── README.md                           # Dokumentasi proyek
```

---

## 2. Penjelasan Rinci File, Syntax, Nama Variabel, dan Fungsi

---

### 2.1. `go.mod` & `go.sum` (Konfigurasi Modul & Dependensi)

* **Fungsi File:** Mendefinisikan nama modul Go (`bab6`) serta mencatat seluruh versi pustaka pihak ketiga (*third-party packages*) yang digunakan proyek.
* **Variabel & Syntax Penting:**
  - `module bab6`: Mendeklarasikan nama modul proyek.
  - `require`: Blok yang memuat dependensi luar:
    - `github.com/gin-gonic/gin`: Framework web HTTP berkinerja tinggi.
    - `github.com/gin-contrib/cors`: Middleware untuk menangani Cross-Origin Resource Sharing.
    - `github.com/lib/pq`: Driver PostgreSQL murni untuk Go.
    - `github.com/rubenv/sql-migrate`: Library migrasi database SQL.
  - `go.sum`: Menyimpan hash SHA-256 untuk memverifikasi integritas versi library yang diunduh.

---

### 2.2. `entities/card_bridge_entity.go` (Model Data)

* **Lokasi File:** `entities/card_bridge_entity.go`
* **Kode:**
  ```go
  package entities

  type Card struct {
  	ID string `json: "id"`
  }
  ```
* **Fungsi File:** Mendefinisikan struktur data (*struct*) representasi entity tabel `card` di database.
* **Penjelasan Syntax & Variabel:**
  - `package entities`: Menandakan file ini berada di dalam package `entities`.
  - `type Card struct`: Mendeklarasikan tipe data bentukan bernama `Card` (huruf kapital di awal agar *exported* / publik).
  - `ID string`: Variabel/field bertipe `string` yang menampung nilai ID kartu.
  - `` `json: "id"` ``: *Struct tag* untuk mengatur nama field saat dikonversi menjadi JSON (misal pada response API).

---

### 2.3. `databases/sql_migrations/1_initiate.sql` (Skrip Migrasi SQL)

* **Lokasi File:** `databases/sql_migrations/1_initiate.sql`
* **Kode:**
  ```sql
  -- +migrate Up
  -- +migrate StatementBegin

  CREATE TABLE card (
      id varchar(20)
  );

  -- +migrate StatementEnd
  ```
* **Fungsi File:** File skrip SQL DDL untuk membuat tabel `card` secara otomatis saat aplikasi melakukan migrasi pertama kali.
* **Penjelasan Syntax:**
  - `-- +migrate Up`: Perintah direktif untuk `sql-migrate` agar mengeksekusi skrip saat migrasi maju.
  - `-- +migrate StatementBegin` / `-- +migrate StatementEnd`: Pembatas blok transaksi SQL.
  - `CREATE TABLE card ( id varchar(20) );`: Perintah SQL untuk membuat tabel `card` dengan satu kolom `id` berjenis `varchar(20)`.

---

### 2.4. `databases/database.go` (Koneksi & Migrasi Otomatis)

* **Lokasi File:** `databases/database.go`
* **Kode:**
  ```go
  package database

  import (
  	"database/sql"
  	"embed"
  	"fmt"

  	migrate "github.com/rubenv/sql-migrate"
  )

  //go:embed sql_migrations/*.sql
  var dbMigration embed.FS
  var DbConnection *sql.DB

  func DBMigration(dbParam *sql.DB) {
  	migrations := &migrate.EmbedFileSystemMigrationSource{
  		FileSystem: dbMigration,
  		Root:       "sql_migrations",
  	}
  	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)

  	if errs != nil {
  		panic(errs)
  	}
  	DbConnection = dbParam
  	fmt.Println("Migrasi berhasil dilakukan", n, migrations)
  }
  ```
* **Fungsi File:** Menyematkan file SQL migrasi ke dalam binary program menggunakan `embed`, menjalankan migrasi ke PostgreSQL, dan menyimpan pointer koneksi database global.
* **Penjelasan Syntax & Variabel:**
  - `package database`: Package manajemen database.
  - `//go:embed sql_migrations/*.sql`: Direktif compiler Go untuk menyertakan file `.sql` ke dalam variabel `dbMigration`.
  - `dbMigration embed.FS`: Variabel sistem file virtual (*embedded file system*) penampung file `.sql`.
  - `DbConnection *sql.DB`: Variabel global pointer koneksi SQL yang dapat diakses oleh package lain.
  - `func DBMigration(dbParam *sql.DB)`: Fungsi untuk mengeksekusi migrasi SQL menggunakan `sql-migrate`. Parameter `dbParam` adalah koneksi database aktif yang diteruskan dari `main.go`.
  - `migrate.Exec(...)`: Menjalankan migrasi database ke arah `migrate.Up`. Variabel `n` menampung jumlah file migrasi yang sukses dijalankan.
  - `panic(errs)`: Menghentikan program jika terjadi error saat migrasi.

---

### 2.5. `repositories/card_bridge_repo.go` (Akses Database / Query SQL)

* **Lokasi File:** `repositories/card_bridge_repo.go`
* **Kode:**
  ```go
  package repositories

  import (
  	"bab6/entities"
  	"database/sql"
  )

  func GetCards(db *sql.DB) (result []entities.Card, err error) {
  	sql := "SELECT * FROM card"
  	rows, err := db.Query(sql)

  	if err != nil {
  		return
  	}
  	defer rows.Close()

  	for rows.Next() {
  		var data entities.Card
  		err = rows.Scan(&data.ID)
  		if err != nil {
  			return
  		}
  		result = append(result, data)
  	}
  	return
  }

  func InsertCard(db *sql.DB, card entities.Card) (err error) {
  	sql := "INSERT INTO card(id) values($1)"
  	_, err = db.Exec(sql, card.ID)
  	return err
  }

  func DeleteCard(db *sql.DB, card entities.Card) (err error) {
  	sql := "DELETE FROM card WHERE id = $1"
  	_, err = db.Exec(sql, card.ID)
  	return err
  }
  ```
* **Fungsi File:** Berada pada *Repository Layer*, bertanggung jawab murni untuk menjalankan query SQL (`SELECT`, `INSERT`, `DELETE`) ke PostgreSQL tanpa melibatkan logika HTTP.
* **Penjelasan Syntax & Variabel:**
  - `GetCards(db *sql.DB) (result []entities.Card, err error)`:
    - Fungsi mengambil seluruh data kartu.
    - Variabel `result`: Slice bertipe `entities.Card` untuk menampung banyak data kartu.
    - Variabel `rows`: Objek `*sql.Rows` hasil eksekusi `db.Query(sql)`.
    - `defer rows.Close()`: Memastikan pointer `rows` ditutup setelah fungsi selesai agar tidak terjadi *connection leak*.
    - `rows.Next()` & `rows.Scan(&data.ID)`: Melakukan iterasi baris data hasil query dan memetakannya ke memori `&data.ID`.
    - `append(result, data)`: Menambahkan data kartu ke dalam slice `result`.
  - `InsertCard(db *sql.DB, card entities.Card) (err error)`:
    - Fungsi menambahkan data baru.
    - Menggunakan parameter placeholder `$1` (`INSERT INTO card(id) values($1)`) untuk mencegah *SQL Injection*.
    - `db.Exec(...)`: Mengeksekusi perintah SQL tanpa mengembalikan baris data (`INSERT`/`UPDATE`/`DELETE`).
  - `DeleteCard(db *sql.DB, card entities.Card) (err error)`:
    - Fungsi menghapus data berdasarkan `id` (`DELETE FROM card WHERE id = $1`).

---

### 2.6. `controllers/card_bridge_controllers.go` (Logika Bisnis & HTTP Handler)

* **Lokasi File:** `controllers/card_bridge_controllers.go`
* **Kode:**
  ```go
  package controllers

  import (
  	database "bab6/databases"
  	"bab6/entities"
  	"bab6/repositories"
  	"net/http"

  	"github.com/gin-gonic/gin"
  )

  func GetCards(c *gin.Context) {
  	var result gin.H

  	card, err := repositories.GetCards(database.DbConnection)

  	if err != nil {
  		result = gin.H{
  			"result": err.Error(),
  		}
  	} else {
  		result = gin.H{
  			"result": card,
  		}
  	}
  	c.JSON(http.StatusOK, result)
  }

  func InsertCard(c *gin.Context) {
  	var card entities.Card
  	idCard := c.Param("id")
  	card.ID = idCard

  	err := repositories.InsertCard(database.DbConnection, card)

  	if err != nil {
  		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  		return
  	}
  	c.JSON(http.StatusOK, card)
  }

  func DeleteCard(c *gin.Context) {
  	var card entities.Card
  	idCard := c.Param("id")
  	card.ID = idCard

  	err := repositories.DeleteCard(database.DbConnection, card)

  	if err != nil {
  		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  		return
  	}
  	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus", "id": idCard})
  }
  ```
* **Fungsi File:** Berada pada *Controller Layer*, bertugas menerima HTTP request dari klien melalui objek `*gin.Context`, memanggil repository, dan mengembalikan response berformat JSON.
* **Penjelasan Syntax & Variabel:**
  - `GetCards(c *gin.Context)`:
    - Variabel `c`: Pointer ke `gin.Context` yang merepresentasikan konteks request dan response HTTP.
    - `var result gin.H`: `gin.H` (`map[string]interface{}`) untuk menyusun struktur JSON dinamis.
    - `c.JSON(http.StatusOK, result)`: Mengirimkan response HTTP `200 OK` dengan payload JSON.
  - `InsertCard(c *gin.Context)`:
    - `idCard := c.Param("id")`: Mengambil parameter dinamis `:id` dari URL.
    - `c.JSON(http.StatusInternalServerError, ...)`: Mengirim response HTTP status `500` jika terjadi kesalahan database.
  - `DeleteCard(c *gin.Context)`:
    - Menghapus kartu berdasarkan parameter `:id` URL dan mengembalikan pesan sukses `"Data Berhasil Dihapus"`.

---

### 2.7. `routers/card_bridge_routers.go` (Routing & CORS Middleware)

* **Lokasi File:** `routers/card_bridge_routers.go`
* **Kode:**
  ```go
  package routers

  import (
  	"bab6/controllers"
  	"time"

  	"github.com/gin-contrib/cors"
  	"github.com/gin-gonic/gin"
  )

  func StartServer() *gin.Engine {
  	router := gin.Default()

  	router.Use(cors.New(cors.Config{
  		AllowOrigins:     []string{"*"},
  		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
  		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
  		ExposeHeaders:    []string{"Content-Length"},
  		AllowCredentials: true,
  		MaxAge:           12 * time.Hour,
  	}))
  	router.GET("/cards", controllers.GetCards)
  	router.POST("/card/input/:id", controllers.InsertCard)
  	router.DELETE("/card/delete/:id", controllers.DeleteCard)

  	return router
  }
  ```
* **Fungsi File:** Menginisialisasi instance router Gin, mengatur konfigurasi CORS (*Cross-Origin Resource Sharing*), dan memetakan endpoint URL ke fungsi controller yang sesuai.
* **Penjelasan Syntax & Variabel:**
  - `StartServer() *gin.Engine`: Fungsi yang mengembalikan pointer instance engine Gin.
  - `router := gin.Default()`: Membuat router Gin dengan middleware bawaan (Logger & Recovery).
  - `router.Use(cors.New(...))`: Mengaktifkan middleware CORS agar API dapat diakses oleh domain/klien luar (misal frontend React/Vue atau Postman).
  - Pendaftaran Endpoint:
    - `router.GET("/cards", controllers.GetCards)`
    - `router.POST("/card/input/:id", controllers.InsertCard)`
    - `router.DELETE("/card/delete/:id", controllers.DeleteCard)`

---

### 2.8. `main.go` (Entry Point / Titik Masuk Utama)

* **Lokasi File:** `main.go`
* **Kode:**
  ```go
  package main

  import (
  	database "bab6/databases"
  	"bab6/routers"
  	"database/sql"
  	"fmt"
  	"log"

  	_ "github.com/lib/pq"
  )

  const (
  	host     = "localhost"
  	port     = 5432
  	user     = "postgres"
  	password = "nickytirta54" // SESUAIKAN DENGAN PASSWORD POSTGRE YANG TELAH DIDAFTARKAN
  	dbName   = "mcsbab6"      // SESUAIKAN DENGAN NAMA DATABASE YANG DIBUAT
  )

  var (
  	DB  *sql.DB
  	err error
  )

  func main() {
  	var PORT = ":50000" // ISI DENGAN PORT BEBAS (DAPAT DICARI DI GOOGLE DENGAN KEYWORD "UNUSED PORT RANGE")

  	psqlInfo := fmt.Sprintf(
  		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
  		host, port, user, password, dbName,
  	)

  	DB, err = sql.Open("postgres", psqlInfo)

  	if err != nil {
  		log.Fatal("Error Open DB: ", err)
  	}

  	database.DBMigration(DB)

  	defer DB.Close()

  	routers.StartServer().Run(PORT)
  	fmt.Printf("Success Connected")
  }
  ```
* **Fungsi File:** File utama tempat program dieksekusi. Berfungsi mengatur koneksi database PostgreSQL, memanggil fungsi migrasi tabel, dan menjalankan server HTTP di port tertentu.
* **Penjelasan Syntax & Variabel:**
  - `package main`: Package utama aplikasi Go *executable*.
  - `_ "github.com/lib/pq"`: *Side-effect import* untuk mendaftarkan driver PostgreSQL ke package standar `database/sql`.
  - Blok `const (...)`: Menyimpan konfigurasi koneksi database:
    - `host`: Alamat server database (`"localhost"`).
    - `port`: Port PostgreSQL (`5432`).
    - `user`: Username database (`"postgres"`).
    - `password`: Password database (`"nickytirta54"`).
    - `dbName`: Nama database (`"mcsbab6"`).
  - `psqlInfo`: Variabel string DSN (*Data Source Name*) yang merangkai konfigurasi koneksi PostgreSQL.
  - `sql.Open("postgres", psqlInfo)`: Membuka koneksi awal ke database.
  - `database.DBMigration(DB)`: Menjalankan migrasi database otomatis.
  - `routers.StartServer().Run(PORT)`: Menjalankan HTTP server pada port `:50000`.

---

## 3. Dokumentasi Endpoint API

Server berjalan di: **`http://localhost:50000`**

### 1. Tambah Data Kartu (`POST`)
- **URL**: `http://localhost:50000/card/input/:id`
- **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:50000/card/input/CARD-001
  ```
- **Response Sukses (`200 OK`)**:
  ```json
  {
    "id": "CARD-001"
  }
  ```

### 2. Lihat Semua Kartu (`GET`)
- **URL**: `http://localhost:50000/cards`
- **Contoh cURL**:
  ```bash
  curl -X GET http://localhost:50000/cards
  ```
- **Response Sukses (`200 OK`)**:
  ```json
  {
    "result": [
      {
        "id": "CARD-001"
      }
    ]
  }
  ```

### 3. Hapus Data Kartu (`DELETE`)
- **URL**: `http://localhost:50000/card/delete/:id`
- **Contoh cURL**:
  ```bash
  curl -X DELETE http://localhost:50000/card/delete/CARD-001
  ```
- **Response Sukses (`200 OK`)**:
  ```json
  {
    "id": "CARD-001",
    "message": "Data Berhasil Dihapus"
  }
  ```

---

## 4. Cara Menjalankan Proyek

1. **Buat Database PostgreSQL**
   ```sql
   CREATE DATABASE mcsbab6;
   ```
2. **Sesuaikan Konfigurasi** di `main.go` (`password` dan `dbName`).
3. **Unduh Dependensi**
   ```bash
   go mod tidy
   ```
4. **Jalankan Aplikasi**
   ```bash
   go run main.go
   ```
