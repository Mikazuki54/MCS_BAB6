# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode Sangat Mendalam, Bedah Simbol/Sintaks (seperti `_`, `*`, `:=`, `&`), Serta Bedah Tuntas Setiap Function di Masing-masing File.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Sebelum masuk ke kode, mari kita pahami apa yang dibuat oleh aplikasi ini dengan analogi sederhana: **Restoran**.

```
[ Pelanggan / Browser / Postman ]
               │
               ▼  1. Pesan/Update Makanan (HTTP Request: GET, POST, PUT, DELETE)
┌──────────────────────────────────────────────┐
│  A. ROUTER (Buku Menu & Kasir)              │ -> Menerima pesanan dan mengarahkan ke Pelayan yang tepat.
└──────────────────────────────────────────────┘
               │
               ▼  2. Meneruskan Pesanan
┌──────────────────────────────────────────────┐
│  B. CONTROLLER (Pelayan / Waiter)           │ -> Memeriksa format pesanan & menyajikan makanan (JSON Response).
└──────────────────────────────────────────────┘
               │
               ▼  3. Minta Bahan / Olah Data
┌──────────────────────────────────────────────┐
│  C. REPOSITORY (Koki / Chef di Dapur)        │ -> Membaca & Menulis data langsung ke Gudang (Database SQL).
└──────────────────────────────────────────────┘
               │
               ▼  4. Ambil / Simpan Data
┌──────────────────────────────────────────────┐
│  D. DATABASE & ENTITY (Gudang & Bentuk Wadah)│ -> PostgreSQL (Gudang Data) & Card Struct (Bentuk Wadah Kartu).
└──────────────────────────────────────────────┘
```

* **REST API**: Cara aplikasi saling berkomunikasi menggunakan protokol HTTP. Klien mengirim *Request* (permintaan) dan server memberikan *Response* (jawaban berupa teks JSON).
* **HTTP Methods**:
  * **`GET`**: Membaca / mengambil data (seperti melihat menu).
  * **`POST`**: Menambahkan data baru (seperti memesan makanan baru).
  * **`PUT`**: Memperbarui/mengubah data yang *sudah ada* (seperti mengganti isi pesanan yang salah).
  * **`DELETE`**: Menghapus data (seperti membatalkan pesanan).
* **Layered Architecture**: Memisah-misah kode ke beberapa folder agar jika ada perubahan atau kerusakan, kita tahu persis folder mana yang harus diperbaiki tanpa mengganggu bagian lain.

---

## 📖 2. Glosarium Simbol & Istilah Khusus di Golang (Wajib Baca untuk Pemula)

Sebelum membaca kode di bawah, pahami dulu arti dari simbol-simbol khusus yang sering muncul dalam bahasa Go:

| Simbol / Istilah | Nama Resmi | Penjelasan Sederhana untuk Pemula |
| :--- | :--- | :--- |
| **`_`** | **Blank Identifier** | Pengganti variabel penampung kosong. Digunakan untuk **mengabaikan** nilai atau hasil kembalian (*return value*) yang tidak ingin kita pakai di dalam kode. |
| **`*`** | **Pointer / Dereference** | Menunjuk ke alamat memori asli suatu data. Contoh `*sql.DB` artinya variabel ini memegang kunci alamat memori asli dari koneksi database. |
| **`&`** | **Address-of Operator** | Mengambil alamat memori dari sebuah variabel. Contoh `&data.ID` artinya kita meminta alamat kotak memori variabel `data.ID` agar bisa diisi langsung oleh database. |
| **`:=`** | **Short Variable Declaration** | Cara cepat membuat variabel baru sekaligus mengisikan nilainya tanpa perlu menulis kata kunci `var`. |
| **`[]`** | **Slice** | Kumpulan data atau daftar (list) berurutan yang ukurannya bisa bertambah atau berkurang (mirip *Array* di bahasa lain). |
| **`defer`** | **Defer Statement** | Perintah tunda yang menyuruh Go mengeksekusi suatu baris kode (misal `rows.Close()`) persis di akhir sebelum fungsi selesai. |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go Baru Menggunakan Framework Gin & Library Pendukung

### Langkah 1: Install Golang
Pastikan Go sudah terinstal di komputer Anda. Cek melalui terminal:
```bash
go version
```

### Langkah 2: Buat Folder Proyek Baru
```bash
mkdir proyek-baru
cd proyek-baru
```

### Langkah 3: Inisialisasi Modul Go
```bash
go mod init proyek-baru
```

### Langkah 4: Install Pustaka / Dependensi Eksternal
1. **Gin Gonic (Web Server)**: `go get -u "github.com/gin-gonic/gin"`
2. **PostgreSQL Driver**: `go get -u "github.com/lib/pq"`
3. **SQL Migrate**: `go get -u "github.com/rubenv/sql-migrate"`
4. **Godotenv**: `go get -u "github.com/joho/godotenv"`

---

## 📘 5. Penjelasan Rinci Kode, Library Import, dan Bedah Setiap Function Secara Menyeluruh

---

### 5.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

**Tugas File:** Mendefinisikan bentuk data kartu yang dipakai dalam program Go dan format teks JSON saat dikirim ke pemanggil API.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Penjelasan Import & Baris demi Baris:
* **Library yang di-import**: *Tidak ada*.
* **Penjelasan Kode**:
  * ◦ `package entities` : Menandai bahwa file ini milik kelompok folder `entities`.
  * ◦ `type Card struct` : Membuat cetakan objek data baru bernama `Card`.
  * ◦ `ID string` : Properti teks di dalam struct `Card`.
  * ◦ `json: "id"` : *Struct tag* untuk mengubah nama field menjadi `"id"` huruf kecil pada format JSON.

---

### 5.2. `databases/sql_migrations/1_initiate.sql` (Skrip Pembuat Tabel)

**Tugas File:** Membuat tabel `card` di PostgreSQL secara otomatis saat aplikasi pertama kali berjalan.

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE card (
    id varchar(20)
);

-- +migrate StatementEnd
```

#### Penjelasan File:
* **Library yang di-import**: *Tidak ada*.
* **Penjelasan Kode**:
  * ◦ `-- +migrate Up` : Perintah direktif migrasi maju.
  * ◦ `CREATE TABLE card ( id varchar(20) );` : Perintah SQL untuk membuat tabel `card`.

---

### 5.3. `databases/database.go` (Pengelola Koneksi DB & Migrasi Otomatis)

**Tugas File:** Membaca file `.sql`, mengeksekusi pembuatan tabel ke PostgreSQL secara otomatis, dan menyimpan koneksi database.

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

#### A. Penjelasan Library / Folder yang Di-import:
* `"database/sql"`: Library standar Go untuk interaksi SQL.
* `"embed"`: Fitur Go untuk menyertakan file fisik `.sql` ke dalam program binary.
* `"fmt"`: Library standar Go untuk mencetak teks ke terminal.
* `migrate "github.com/rubenv/sql-migrate"`: Library migrasi SQL otomatis.

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `DBMigration(dbParam *sql.DB)`**
  * **Tujuan**: Menjalankan migrasi database otomatis.
  * **Parameter**: `dbParam *sql.DB` (Pointer koneksi database aktif yang dikirim dari `main.go`).
  * **Bedah Isi Fungsi Baris demi Baris**:
    1. `migrations := &migrate.EmbedFileSystemMigrationSource{ ... }`: Menyiapkan sumber file migrasi dari sistem file virtual `dbMigration` di folder `sql_migrations`.
    2. `n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)`: Mengeksekusi migrasi.
       * `n`: Angka jumlah file migrasi yang sukses dijalankan.
       * `errs`: Wadah penampung error jika migrasi gagal.
    3. `if errs != nil { panic(errs) }`: Jika terjadi error, program dihentikan paksa dengan `panic()`.
    4. `DbConnection = dbParam`: Menyimpan koneksi database ke variabel global `DbConnection`.
    5. `fmt.Println(...)`: Mencetak pesan sukses ke terminal.

---

### 5.4. `repositories/card_bridge_repo.go` (Akses Data / Query SQL)

**Tugas File:** Tempat koki menjalankan perintah query SQL langsung (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) ke database PostgreSQL.

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

func UpdateCard(db *sql.DB, card entities.Card) (err error) {
	sql := "UPDATE card SET id = $1 WHERE id = $2"
	_, err = db.Exec(sql, card.ID, "id_lama")
	return err
}

func DeleteCard(db *sql.DB, card entities.Card) (err error) {
	sql := "DELETE FROM card WHERE id = $1"
	_, err = db.Exec(sql, card.ID)
	return err
}
```

#### A. Penjelasan Library / Folder yang Di-import:
* `"bab6/entities"`: Memanggil struct `Card`.
* `"database/sql"`: Memanggil pustaka standar tipe data SQL.

#### B. Bedah Setiap Fungsi (Function) di File Ini:

1. **Fungsi: `GetCards(db *sql.DB) (result []entities.Card, err error)`**
   * **Tujuan**: Mengambil seluruh daftar data kartu dari database.
   * **Bedah Baris demi Baris**:
     * `sql := "SELECT * FROM card"`: Menyiapkan query string SQL.
     * `rows, err := db.Query(sql)`: Menjalankan query `SELECT`. Variabel `rows` menampung banyak baris hasil database.
     * `defer rows.Close()`: Menunda penutupan kursor baris hingga fungsi selesai.
     * `for rows.Next()`: Perulangan membaca baris demi baris.
       * `var data entities.Card`: Wadah penampung sementara 1 baris kartu.
       * `err = rows.Scan(&data.ID)`: **`&` (Address-of)** mengambil alamat memori `data.ID` agar database bisa mengisi langsung nilainya.
       * `result = append(result, data)`: Memasukkan data ke slice `result`.

2. **Fungsi: `InsertCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Menambahkan data kartu baru ke database.
   * **Bedah Baris demi Baris**:
     * `sql := "INSERT INTO card(id) values($1)"`: Query insert dengan parameter `$1`.
     * `_, err = db.Exec(sql, card.ID)`: Menjalankan eksekusi query.
       * **Penjelasan Simbol `_` (Blank Identifier)**: Fungsi `db.Exec` mengembalikan dua nilai yaitu `sql.Result` (metadata hasil eksekusi) dan `error`. Karena kita **tidak butuh** metadata `sql.Result` tersebut, kita wajib menuliskan tanda **`_`** sebagai pengganti variabel untuk **mengabaikannya**. Di Go, variabel yang sudah dibuat wajib dipakai; jika tidak dipakai, compiler akan protes error. Tanda `_` adalah solusi legal untuk mengabaikannya.
     * `return err`: Mengembalikan status error.

3. **Fungsi: `UpdateCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Memperbarui data kartu.
   * **Bedah Baris demi Baris**:
     * `sql := "UPDATE card SET id = $1 WHERE id = $2"`: Query SQL untuk update data.
     * `_, err = db.Exec(sql, card.ID, "id_lama")`: Menggunakan tanda **`_`** untuk **mengabaikan** nilai `sql.Result` yang tidak terpakai, dan hanya mengambil `err`.
     * `return err`: Mengembalikan status error.

4. **Fungsi: `DeleteCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Menghapus data kartu berdasarkan ID.
   * **Bedah Baris demi Baris**:
     * `sql := "DELETE FROM card WHERE id = $1"`: Query SQL hapus data.
     * `_, err = db.Exec(sql, card.ID)`: Menggunakan tanda **`_`** untuk **mengabaikan** return value pertama (`sql.Result`), karena kita hanya peduli pada status `err`.
     * `return err`: Mengembalikan status error.

---

### 5.6. `controllers/card_bridge_controllers.go` (Logika Bisnis & HTTP Handler)

**Tugas File:** Pelayan yang menerima permintaan HTTP, memanggil repository, dan mengirim balasan JSON.

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

func UpdateCard(c *gin.Context) {
	var card entities.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repositories.UpdateCard(database.DbConnection, card)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diupdate"})
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

#### A. Penjelasan Library / Folder yang Di-import:
* `database "bab6/databases"`: Meng-import package database lokal untuk mengakses koneksi.
* `"bab6/entities"`: Meng-import struct `Card`.
* `"bab6/repositories"`: Meng-import fungsi repository SQL.
* `"net/http"`: Library HTTP status code.
* `"github.com/gin-gonic/gin"`: Framework Gin untuk manajemen request/response.

#### B. Bedah Setiap Fungsi (Function) di File Ini:

1. **Fungsi: `GetCards(c *gin.Context)`**
   * **Tujuan**: Menangani request `GET /cards`.
   * **Bedah Baris demi Baris**:
     * `var result gin.H`: Membuat map dinamis JSON.
     * `card, err := repositories.GetCards(database.DbConnection)`: Memanggil repository untuk mengambil data.
     * `if err != nil`: Percabangan untuk mengisi `result` dengan pesan error atau data kartu.
     * `c.JSON(http.StatusOK, result)`: Mengirim response JSON status `200 OK`.

2. **Fungsi: `InsertCard(c *gin.Context)`**
   * **Tujuan**: Menangani request `POST /card/input/:id`.
   * **Bedah Baris demi Baris**:
     * `idCard := c.Param("id")`: Mengambil parameter `:id` dari URL.
     * `err := repositories.InsertCard(...)`: Menyimpan data ke database.
     * `if err != nil`: Jika gagal, kirim status `500` (`http.StatusInternalServerError`).
     * `c.JSON(http.StatusOK, card)`: Kirim status `200 OK` jika sukses.

3. **Fungsi: `UpdateCard(c *gin.Context)`**
   * **Tujuan**: Menangani request `PUT /card/update/:id`.
   * **Bedah Baris demi Baris**:
     * `if err := c.ShouldBindJSON(&card); err != nil`: Membaca body JSON dari klien. Jika gagal, kirim status `400 Bad Request`.
     * `c.JSON(http.StatusOK, ...)`: Kirim pesan sukses update.

4. **Fungsi: `DeleteCard(c *gin.Context)`**
   * **Tujuan**: Menangani request `DELETE /card/delete/:id`.
   * **Bedah Baris demi Baris**:
     * Mengambil ID dari URL, menghapus data via repository, lalu mengirim response sukses `200 OK`.

---

### 5.7. `routers/card_bridge_routers.go` (Pendaftaran Alamat URL & CORS)

**Tugas File:** Mengatur daftar alamat URL (endpoint), menentukan HTTP Method, dan mengatur izin CORS.

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
	router.PUT("/card/update/:id", controllers.UpdateCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### A. Penjelasan Library / Folder yang Di-import:
* `"bab6/controllers"`: Memanggil controllers.
* `"time"`: Pengaturan durasi waktu cache CORS.
* `"github.com/gin-contrib/cors"`: Middleware CORS resmi Gin.
* `"github.com/gin-gonic/gin"`: Framework Gin.

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `StartServer() *gin.Engine`**
  * **Tujuan**: Menginisialisasi router Gin dan mendaftarkan endpoint rute API.
  * **Bedah Baris demi Baris**:
    1. `router := gin.Default()`: Membuat mesin router Gin baru.
    2. `router.Use(cors.New(...))`: Mengatur izin CORS.
    3. `router.GET`, `router.POST`, `router.PUT`, `router.DELETE`: Mendaftarkan endpoint URL ke fungsi controller.
    4. `return router`: Mengembalikan objek router.

---

### 5.8. `main.go` (Pintu Masuk Utama Aplikasi)

**Tugas File:** Titik awal eksekusi program. Menyiapkan parameter database, membuka koneksi PostgreSQL, memanggil migrasi, dan menyalakan server web.

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
	password = "nickytirta54"
	dbName   = "mcsbab6"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	var PORT = ":50000"

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

#### A. Penjelasan Library / Folder yang Di-import:
* `database "bab6/databases"`: Package migrasi database lokal.
* `"bab6/routers"`: Package routers server.
* `"database/sql"`: Library interaksi SQL Go.
* `"fmt"`: Pemformatan string.
* `"log"`: Logging error fatal.
* `_ "github.com/lib/pq"`: **Side-Effect Import** (menggunakan tanda **`_`** di depan untuk mendaftarkan driver PostgreSQL ke package standar `database/sql` tanpa memanggil fungsinya secara langsung).

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `main()`**
  * **Tujuan**: Fungsi utama eksekusi program.
  * **Bedah Baris demi Baris**:
    1. `var PORT = ":50000"`: Menentukan port server.
    2. `psqlInfo := fmt.Sprintf(...)`: Merangkai koneksi string database.
    3. `DB, err = sql.Open(...)`: Membuka koneksi database.
    4. `database.DBMigration(DB)`: Menjalankan migrasi database.
    5. `defer DB.Close()`: Menutup koneksi database di akhir program.
    6. `routers.StartServer().Run(PORT)`: Menyalakan web server.

---

## 🚀 6. Cara Menjalankan Proyek Utama (`bab6`)

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

---

## 🧪 7. Pengujian Endpoint API (Termasuk PUT)

Base URL Server: **`http://localhost:50000`**

### 1. Tambah Kartu (`POST`)
* **URL**: `http://localhost:50000/card/input/CARD-001`
* **cURL**:
  ```bash
  curl -X POST http://localhost:50000/card/input/CARD-001
  ```

### 2. Lihat Semua Kartu (`GET`)
* **URL**: `http://localhost:50000/cards`
* **cURL**:
  ```bash
  curl -X GET http://localhost:50000/cards
  ```

### 3. Perbarui Kartu (`PUT`)
* **URL**: `http://localhost:50000/card/update/CARD-001`
* **cURL**:
  ```bash
  curl -X PUT http://localhost:50000/card/update/CARD-001 \
  -H "Content-Type: application/json" \
  -d '{"id": "CARD-BARU"}'
  ```

### 4. Hapus Kartu (`DELETE`)
* **URL**: `http://localhost:50000/card/delete/CARD-001`
* **cURL**:
  ```bash
  curl -X DELETE http://localhost:50000/card/delete/CARD-001
  ```
