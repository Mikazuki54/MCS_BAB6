# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode 100% Tanpa Ada Baris yang Terlewat, Termasuk Penjelasan Simbol, Compiler Directive (`//go:embed`), Side-Effect Import (`_`), dan Bedah Setiap Function.*

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
  * **`POST`**: Menciptakan/menambahkan data baru (seperti memesan makanan baru).
  * **`PUT`**: Memperbarui/mengubah data yang *sudah ada* (seperti mengganti isi pesanan yang salah).
  * **`DELETE`**: Menghapus data (seperti membatalkan pesanan).
* **Layered Architecture**: Memisah-misah kode ke beberapa folder agar jika ada perubahan atau kerusakan, kita tahu persis folder mana yang harus diperbaiki tanpa mengganggu bagian lain.

---

## 📖 2. Glosarium Simbol, Compiler Directive, CORS, & Istilah Khusus di Golang

Sebelum membaca kode di bawah, pahami dulu arti dari simbol-simbol khusus dan istilah penting:

| Istilah / Simbol | Nama Resmi | Penjelasan Sederhana untuk Pemula |
| :--- | :--- | :--- |
| **`//go:embed`** | **Compiler Directive** | Perintah khusus untuk menyertakan file atau folder fisik (seperti file `.sql`) langsung ke dalam file program utama saat di-compile, sehingga file eksternal tidak perlu ikut disebar. |
| **CORS** | **Cross-Origin Resource Sharing** | Satpam browser yang mengatur apakah website dari domain lain boleh mengakses API kita. |
| **`_`** | **Blank Identifier** | Digunakan untuk **mengabaikan** nilai hasil kembalian (*return value*) fungsi yang tidak kita gunakan (wajib di Go agar tidak error). |
| **`*`** | **Pointer** | Tanda yang menunjuk langsung ke alamat memori asli data (contoh `*sql.DB`). |
| **`&`** | **Address-of Operator** | Tanda untuk mengambil alamat memori variabel agar nilainya bisa diisi langsung oleh database. |
| **`:=`** | **Short Variable Declaration** | Cara cepat membuat dan mengisi variabel baru tanpa kata kunci `var`. |
| **`defer`** | **Defer Statement** | Perintah untuk menunda eksekusi suatu baris kode sampai fungsi utama selesai. |

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
2. **CORS Middleware**: `go get -u "github.com/gin-contrib/cors"`
3. **PostgreSQL Driver**: `go get -u "github.com/lib/pq"`
4. **SQL Migrate**: `go get -u "github.com/rubenv/sql-migrate"`

---

## 📘 5. Bedah Tuntas 100% Baris demi Baris Setiap File dan Fungsinya (Tanpa Ada yang Terlewat)

---

### 5.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Bedah Baris demi Baris:
1. `package entities` : Mendeklarasikan bahwa file ini tergabung dalam kelompok folder `entities`.
2. `type Card struct {` : Memulai pembuatan cetakan data baru (*struct*) bernama `Card` (huruf C kapital agar publik/bisa diakses luar package).
3. `ID string` : Membuat kotak atribut bernama `ID` dengan tipe data teks (`string`).
4. `json: "id"` : *Struct tag* penanda agar saat struct ini diubah menjadi JSON, variabel `ID` berubah nama menjadi `"id"` (huruf kecil).
5. `}` : Menutup blok deklarasi struct `Card`.

---

### 5.2. `databases/sql_migrations/1_initiate.sql` (Skrip Pembuat Tabel)

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE card (
    id varchar(20)
);

-- +migrate StatementEnd
```

#### Bedah Baris demi Baris:
1. `-- +migrate Up` : Perintah direktif migrasi maju untuk `sql-migrate`.
2. `-- +migrate StatementBegin` : Menandai awal blok transaksi SQL.
3. `CREATE TABLE card (` : Perintah SQL untuk membuat tabel baru bernama `card`.
4. `id varchar(20)` : Membuat kolom `id` bertipe teks dinamis maksimal 20 karakter.
5. `);` : Menutup perintah pembuatan tabel.
6. `-- +migrate StatementEnd` : Menandai akhir blok transaksi SQL.

---

### 5.3. `databases/database.go` (Koneksi & Migrasi Otomatis)

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

#### Bedah Import & Setiap Baris Kode / Fungsi:
* **Import**:
  * `"database/sql"`: Library standar Go untuk interaksi SQL.
  * `"embed"`: Library Go untuk menyertakan file fisik ke dalam binary.
  * `"fmt"`: Library Go untuk cetak teks ke terminal.
  * `migrate "github.com/rubenv/sql-migrate"`: Library migrasi SQL (diberi alias `migrate`).
* **Variabel Global**:
  * `//go:embed sql_migrations/*.sql`: **Compiler Directive**. Perintah khusus kepada compiler Go untuk mengambil semua file berformat `.sql` di dalam folder `sql_migrations/` dan membungkusnya ke dalam program.
  * `var dbMigration embed.FS`: Variabel sistem file virtual penampung file `.sql` hasil embed.
  * `var DbConnection *sql.DB`: Variabel global penyimpan pointer koneksi database yang bisa dipanggil file lain.
* **Bedah Fungsi: `DBMigration(dbParam *sql.DB)`**
  * `dbParam *sql.DB`: Parameter koneksi database aktif yang dikirim dari `main.go`.
  * `migrations := &migrate.EmbedFileSystemMigrationSource{...}`: Mengatur sumber file migrasi dari sistem file virtual `dbMigration` dengan akar folder `"sql_migrations"`.
  * `n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)`: Menjalankan migrasi ke PostgreSQL. `n` menampung jumlah file migrasi sukses, `errs` menampung error jika gagal.
  * `if errs != nil { panic(errs) }`: Jika migrasi error, hentikan program dengan `panic()`.
  * `DbConnection = dbParam`: Menyimpan koneksi database ke variabel global.
  * `fmt.Println(...)`: Mencetak info sukses migrasi ke terminal.

---

### 5.4. `repositories/card_bridge_repo.go` (Akses Data / Query SQL)

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

#### A. Bedah Import:
* `"bab6/entities"`: Memuat package entities untuk struct `Card`.
* `"database/sql"`: Memuat pustaka SQL standar Go.

#### B. Bedah Setiap Fungsi Baris demi Baris:

1. **Fungsi: `GetCards(db *sql.DB) (result []entities.Card, err error)`**
   * `db *sql.DB`: Parameter koneksi database.
   * `(result []entities.Card, err error)`: *Named return values* (slice penampung data kartu dan status error).
   * `sql := "SELECT * FROM card"`: Menyiapkan query string SQL.
   * `rows, err := db.Query(sql)`: Menjalankan query `SELECT`.
   * `if err != nil { return }`: Keluar fungsi jika query gagal.
   * `defer rows.Close()`: Menunda penutupan `rows` hingga fungsi selesai agar memori aman.
   * `for rows.Next()`: Perulangan membaca baris data satu per satu.
     * `var data entities.Card`: Wadah 1 kartu.
     * `err = rows.Scan(&data.ID)`: Menyalin data kolom database ke memori `&data.ID`.
     * `result = append(result, data)`: Memasukkan kartu ke slice `result`.

2. **Fungsi: `InsertCard(db *sql.DB, card entities.Card) (err error)`**
   * `sql := "INSERT INTO card(id) values($1)"`: Query insert data dengan parameter `$1`.
   * `_, err = db.Exec(sql, card.ID)`: Menjalankan eksekusi query.
     * **Penjelasan Simbol `_` (Blank Identifier)**: `db.Exec` mengembalikan dua nilai (`sql.Result` dan `error`). Karena kita **tidak butuh** nilai pertama (`sql.Result`), kita wajib menuliskan tanda **`_`** untuk **mengabaikannya** agar compiler Go tidak komplain error.
   * `return err`: Mengembalikan status error.

3. **Fungsi: `UpdateCard(db *sql.DB, card entities.Card) (err error)`**
   * `sql := "UPDATE card SET id = $1 WHERE id = $2"`: Query update data.
   * `_, err = db.Exec(...)`: Menggunakan tanda **`_`** untuk **mengabaikan** nilai return pertama (`sql.Result`) yang tidak terpakai.
   * `return err`: Mengembalikan error.

4. **Fungsi: `DeleteCard(db *sql.DB, card entities.Card) (err error)`**
   * `sql := "DELETE FROM card WHERE id = $1"`: Query delete data.
   * `_, err = db.Exec(...)`: Menggunakan tanda **`_`** untuk **mengabaikan** return value pertama.
   * `return err`: Mengembalikan error.

---

### 5.5. `controllers/card_bridge_controllers.go` (HTTP Handler & Logika)

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

#### A. Bedah Import:
* `database "bab6/databases"`: Meng-import package databases lokal.
* `"bab6/entities"`: Meng-import struct `Card`.
* `"bab6/repositories"`: Meng-import repository fungsi SQL.
* `"net/http"`: Library HTTP status code.
* `"github.com/gin-gonic/gin"`: Framework Gin untuk manajemen request/response.

#### B. Bedah Setiap Fungsi Baris demi Baris:

1. **Fungsi: `GetCards(c *gin.Context)`**
   * `c *gin.Context`: Konteks request HTTP.
   * `var result gin.H`: Map JSON dinamis.
   * `card, err := repositories.GetCards(database.DbConnection)`: Panggil repo.
   * `c.JSON(http.StatusOK, result)`: Kirim response JSON `200 OK`.

2. **Fungsi: `InsertCard(c *gin.Context)`**
   * `idCard := c.Param("id")`: Ambil parameter `:id` dari URL.
   * `err := repositories.InsertCard(...)`: Simpan ke DB.
   * `c.JSON(http.StatusInternalServerError, ...)`: Kirim status `500` jika gagal.

3. **Fungsi: `UpdateCard(c *gin.Context)`**
   * `c.ShouldBindJSON(&card)`: Membaca body JSON dari klien.
   * `c.JSON(http.StatusOK, ...)`: Kirim pesan sukses update.

4. **Fungsi: `DeleteCard(c *gin.Context)`**
   * Mengambil ID dari URL, hapus data via repo, kirim response sukses `200 OK`.

---

### 5.6. `routers/card_bridge_routers.go` (Routing & CORS)

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

#### A. Bedah Import:
* `"bab6/controllers"`: Memanggil controllers.
* `"time"`: Pengaturan durasi waktu cache CORS.
* `"github.com/gin-contrib/cors"`: Middleware CORS resmi Gin.
* `"github.com/gin-gonic/gin"`: Framework Gin.

#### B. Bedah Fungsi: `StartServer() *gin.Engine`
* `router := gin.Default()`: Membuat mesin router Gin.
* `router.Use(cors.New(...))`: Mengatur izin CORS.
* `router.GET`, `router.POST`, `router.PUT`, `router.DELETE`: Mendaftarkan endpoint URL ke controller.
* `return router`: Mengembalikan objek router.

---

### 5.7. `main.go` (Pintu Masuk Utama Aplikasi)

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

#### A. Bedah Import:
* `database "bab6/databases"`: Package migrasi database lokal.
* `"bab6/routers"`: Package routers server.
* `"database/sql"`: Library interaksi SQL Go.
* `"fmt"`: Pemformatan string.
* `"log"`: Logging error fatal.
* `_ "github.com/lib/pq"`: **Side-Effect Import** (menggunakan tanda **`_`** di depan untuk mendaftarkan driver PostgreSQL ke package standar `database/sql` tanpa memanggil fungsinya secara langsung).

#### B. Bedah Fungsi: `main()`
* `var PORT = ":50000"`: Menentukan port server.
* `psqlInfo := fmt.Sprintf(...)`: Merangkai koneksi string database.
* `DB, err = sql.Open(...)`: Membuka koneksi database.
* `database.DBMigration(DB)`: Menjalankan migrasi database.
* `defer DB.Close()`: Menutup koneksi database di akhir program.
* `routers.StartServer().Run(PORT)`: Menyalakan web server.

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
