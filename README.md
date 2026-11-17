# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode Sangat Mendalam, Baris demi Baris, Serta Bedah Tuntas Setiap Function (Fungsi) yang Ada di Masing-masing File.*

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

## 📖 2. Glosarium Istilah Pemula

*Baca ini agar lebih mudah memahami istilah teknis di bawah:*

| Istilah | Penjelasan Sederhana untuk Pemula |
| :--- | :--- |
| **Golang / Go** | Bahasa pemrograman yang cepat, ringan, dan cocok untuk backend server. |
| **Gin Gonic** | Alat bantu (framework) untuk membuat web server di Go. |
| **Function (Fungsi)** | Blok kode mandiri yang diberi nama untuk menjalankan tugas tertentu saat dipanggil. |
| **Struct** | Fitur Go untuk membuat wadah tipe data sendiri (seperti objek). |
| **Pointer (`*`)** | Tanda penunjuk alamat memori komputer asli, bukan menduplikat datanya. |
| **Parameter** | Nilai/bahan masukan yang dikirim ke dalam sebuah fungsi agar bisa diolah. |
| **JSON** | Format teks pertukaran data di internet (`"key": "value"`). |
| **CORS** | Sistem keamanan browser agar web lain tidak sembarangan mengambil data API kita. |
| **`defer`** | Perintah "nanti saja" untuk menjalankan fungsi tepat sebelum fungsi utama selesai. |
| **`import`** | Perintah untuk memanggil atau "membawa masuk" folder/library luar agar bisa digunakan di file kita. |

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

## 📘 5. Penjelasan Rinci Kode, Library Import, dan Bedah Setiap Function Berdasarkan Urutan Alur

---

### 5.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

**Tugas File:** Mendefinisikan bentuk data kartu yang dipakai dalam program Go dan format teks JSON saat dikirim ke pemanggil API.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### A. Penjelasan Import & File:
* **Library yang di-import**: *Tidak ada*, karena file ini murni mendefinisikan struktur data dasar.
* **Penjelasan Kode**:
  * ◦ `package entities` : Menandai bahwa file ini adalah bagian dari kelompok folder `entities`.
  * ◦ `type Card struct` : Membuat cetakan objek data baru bernama `Card`. Karena diawali huruf kapital (`C`), struct ini bisa diakses dari luar folder.
  * ◦ `ID string` : Field/atribut di dalam struct berupa teks (`string`).
  * ◦ `json: "id"` : *Struct tag* untuk mengubah nama field `ID` menjadi `"id"` huruf kecil saat diubah ke format JSON.

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
* **Library yang di-import**: *Tidak ada* (berisi perintah DDL SQL murni).
* **Penjelasan Kode**:
  * ◦ `-- +migrate Up` : Perintah direktif untuk library migrasi database.
  * ◦ `CREATE TABLE card ( id varchar(20) );` : Perintah SQL untuk membuat tabel baru bernama `card` dengan kolom `id` bertipe teks maksimal 20 karakter.

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
* `"database/sql"`: Library standar Go untuk berinteraksi dengan database relational.
* `"embed"`: Fitur bawaan Go untuk menyertakan file fisik (file `.sql`) langsung ke dalam program binary.
* `"fmt"`: Library standar Go untuk memformat dan mencetak teks ke terminal.
* `migrate "github.com/rubenv/sql-migrate"`: Library migrasi database SQL otomatis.

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `DBMigration(dbParam *sql.DB)`**
  * **Tujuan**: Menjalankan skrip migrasi SQL secara otomatis ke database PostgreSQL.
  * **Parameter**:
    * `dbParam *sql.DB`: Kunci koneksi database aktif yang dikirim dari file `main.go`. Tanpa parameter ini, fungsi migrasi tidak tahu ke database mana perintah SQL harus dijalankan.
  * **Bedah Isi Fungsi Baris demi Baris**:
    1. `migrations := &migrate.EmbedFileSystemMigrationSource{ FileSystem: dbMigration, Root: "sql_migrations" }`: Menyiapkan sumber file migrasi dengan menunjuk ke sistem file virtual `dbMigration` di dalam folder `sql_migrations`.
    2. `n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)`: Mengeksekusi migrasi.
       * `n`: Menyimpan angka jumlah file migrasi yang sukses dijalankan.
       * `errs`: Menampung objek error jika migrasi gagal.
       * `"postgres"`: Menentukan jenis database yang dituju.
       * `migrate.Up`: Arah migrasi untuk membuat/menambah struktur tabel.
    3. `if errs != nil { panic(errs) }`: Pemeriksa error. Jika `errs` tidak kosong (artinya gagal), program dihentikan paksa dengan `panic()`.
    4. `DbConnection = dbParam`: Menyimpan koneksi database yang valid ke variabel global `DbConnection`.
    5. `fmt.Println(...)` : Mencetak pesan konfirmasi sukses migrasi ke terminal.

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
* `"bab6/entities"`: Memanggil package `entities` agar bisa memakai cetakan struct `entities.Card`.
* `"database/sql"`: Memanggil pustaka standar Go untuk mengelola tipe data koneksi database SQL (`*sql.DB`).

#### B. Bedah Fungsi-Fungsi (Functions) di File Ini:

1. **Fungsi: `GetCards(db *sql.DB) (result []entities.Card, err error)`**
   * **Tujuan**: Mengambil seluruh daftar data kartu yang ada di dalam database.
   * **Parameter**: `db *sql.DB` (koneksi database).
   * **Return Values**: `result []entities.Card` (kumpulan data kartu) dan `err error` (status error). Menggunakan fitur *Named Return*.
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `sql := "SELECT * FROM card"`: Menyiapkan string query SQL untuk memilih semua data dari tabel `card`.
     * `rows, err := db.Query(sql)`: Menjalankan query `SELECT`. Variabel `rows` menampung banyak baris hasil database.
     * `if err != nil { return }`: Jika query gagal, langsung keluar dari fungsi.
     * `defer rows.Close()`: Menjadwalkan penutupan kursor `rows` di akhir fungsi agar RAM komputer tidak bocor.
     * `for rows.Next()`: Perulangan untuk membaca hasil query baris demi baris.
       * `var data entities.Card`: Wadah penampung sementara untuk 1 baris data kartu.
       * `err = rows.Scan(&data.ID)`: Menyalin nilai kolom dari database ke alamat memori `&data.ID`.
       * `result = append(result, data)`: Memasukkan data kartu ke dalam slice `result`.
     * `return`: Mengembalikan hasil `result` dan `err`.

2. **Fungsi: `InsertCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Menyimpan/menambahkan data kartu baru ke database.
   * **Parameter**: `db *sql.DB` dan `card entities.Card`.
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `sql := "INSERT INTO card(id) values($1)"`: Menyusun perintah SQL `INSERT` dengan parameter pengaman `$1`.
     * `_, err = db.Exec(sql, card.ID)`: Menjalankan perintah SQL perubahan data (`INSERT`). Hasil metadata diabaikan dengan `_`, dan error ditampung di `err`.
     * `return err`: Mengembalikan status error (bernilai `nil` jika sukses).

3. **Fungsi: `UpdateCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Memperbarui data kartu yang sudah ada di database.
   * **Bedah Isi Fungsi**:
     * `sql := "UPDATE card SET id = $1 WHERE id = $2"`: Perintah SQL untuk memperbarui kolom `id` berdasarkan kondisi tertentu.
     * `_, err = db.Exec(sql, card.ID, "id_lama")`: Menjalankan query update dengan mengirimkan parameter data baru dan lama.
     * `return err`: Mengembalikan status error.

4. **Fungsi: `DeleteCard(db *sql.DB, card entities.Card) (err error)`**
   * **Tujuan**: Menghapus data kartu berdasarkan ID.
   * **Bedah Isi Fungsi**:
     * `sql := "DELETE FROM card WHERE id = $1"`: Perintah SQL penghapusan baris data.
     * `_, err = db.Exec(sql, card.ID)`: Menjalankan query `DELETE` ke database.
     * `return err`: Mengembalikan status error.

---

### 5.5. `controllers/card_bridge_controllers.go` (Logika Bisnis & HTTP Handler)

**Tugas File:** Pelayan yang menerima permintaan HTTP, meminta repository memproses ke database, dan mengantarkan balasan berformat JSON.

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
* `database "bab6/databases"`: Meng-import package `databases` buatan sendiri untuk mengakses `database.DbConnection`.
* `"bab6/entities"`: Meng-import struct `Card`.
* `"bab6/repositories"`: Meng-import fungsi repository (*GetCards, InsertCard, dll*).
* `"net/http"`: Library bawaan Go untuk kode status HTTP (`http.StatusOK`, `http.StatusInternalServerError`, dll).
* `"github.com/gin-gonic/gin"`: Framework Gin untuk manajemen request HTTP (`gin.Context`) dan format JSON (`gin.H`).

#### B. Bedah Fungsi-Fungsi (Functions) di File Ini:

1. **Fungsi: `GetCards(c *gin.Context)`**
   * **Tujuan**: Menangani request HTTP `GET /cards` untuk menampilkan seluruh data kartu.
   * **Parameter**: `c *gin.Context` (konteks request dan response Gin).
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `var result gin.H`: Membuat variabel map dinamis (`gin.H`) untuk menyusun struktur JSON.
     * `card, err := repositories.GetCards(database.DbConnection)`: Memanggil fungsi repository untuk mengambil data kartu dari PostgreSQL menggunakan koneksi global.
     * `if err != nil`: Jika gagal, `result` diisi pesan error. Jika sukses, `result` diisi data kartu `card`.
     * `c.JSON(http.StatusOK, result)`: Mengirimkan response HTTP berstatus `200 OK` dengan isi payload berupa JSON dari variabel `result`.

2. **Fungsi: `InsertCard(c *gin.Context)`**
   * **Tujuan**: Menangani request HTTP `POST /card/input/:id` untuk menambah kartu baru.
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `var card entities.Card`: Membuat instance struct `Card` kosong.
     * `idCard := c.Param("id")`: Mengambil nilai parameter dinamis `:id` dari URL.
     * `card.ID = idCard`: Memasukkan ID tersebut ke struct `card`.
     * `err := repositories.InsertCard(database.DbConnection, card)`: Menyimpan data ke database.
     * `if err != nil`: Jika gagal, kirim status `500 Internal Server Error` (`http.StatusInternalServerError`) dan hentikan fungsi dengan `return`.
     * `c.JSON(http.StatusOK, card)`: Jika sukses, kirim status `200 OK` beserta data kartu.

3. **Fungsi: `UpdateCard(c *gin.Context)`**
   * **Tujuan**: Menangani request HTTP `PUT /card/update/:id` untuk memperbarui data kartu.
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `if err := c.ShouldBindJSON(&card); err != nil`: Membaca format JSON yang dikirim pengguna di body request dan menyalinnya ke struct `card`. Jika formatnya salah, kirim status `400 Bad Request`.
     * `err := repositories.UpdateCard(...)`: Menjalankan fungsi update di repository.
     * `c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diupdate"})`: Mengirim pesan sukses jika berhasil.

4. **Fungsi: `DeleteCard(c *gin.Context)`**
   * **Tujuan**: Menangani request HTTP `DELETE /card/delete/:id` untuk menghapus kartu.
   * **Bedah Isi Fungsi Baris demi Baris**:
     * `idCard := c.Param("id")`: Mengambil parameter `:id` dari path URL.
     * `err := repositories.DeleteCard(...)`: Menghapus data dari database via repository.
     * `c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus", "id": idCard})`: Mengirim response konfirmasi sukses penghapusan.

---

### 5.6. `routers/card_bridge_routers.go` (Pendaftaran Alamat URL & CORS)

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
* `"bab6/controllers"`: Memanggil package controllers.
* `"time"`: Library waktu Go untuk pengaturan kedaluwarsa cache CORS.
* `"github.com/gin-contrib/cors"`: Middleware CORS resmi Gin.
* `"github.com/gin-gonic/gin"`: Framework Gin.

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `StartServer() *gin.Engine`**
  * **Tujuan**: Menginisialisasi mesin router Gin, memasang middleware CORS, mendaftarkan semua endpoint URL API, dan mengembalikan objek engine router.
  * **Bedah Isi Fungsi Baris demi Baris**:
    1. `router := gin.Default()`: Membuat engine router Gin baru dengan fitur bawaan *Logger* dan *Recovery*.
    2. `router.Use(cors.New(...))`: Menerapkan aturan CORS agar API bisa diakses oleh aplikasi luar/frontend web tanpa batasan domain.
    3. `router.GET("/cards", controllers.GetCards)`: Mendaftarkan rute `GET` untuk melihat semua kartu.
    4. `router.POST("/card/input/:id", controllers.InsertCard)`: Mendaftarkan rute `POST` untuk menambah kartu baru.
    5. `router.PUT("/card/update/:id", controllers.UpdateCard)`: Mendaftarkan rute `PUT` untuk memperbarui kartu.
    6. `router.DELETE("/card/delete/:id", controllers.DeleteCard)`: Mendaftarkan rute `DELETE` untuk menghapus kartu.
    7. `return router`: Mengembalikan objek `router` yang sudah siap dijalankan di `main.go`.

---

### 5.7. `main.go` (Pintu Masuk Utama Aplikasi)

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
* `database "bab6/databases"`: Meng-import package database lokal untuk migrasi.
* `"bab6/routers"`: Meng-import package routers untuk menjalankan `StartServer()`.
* `"database/sql"`: Library standar Go untuk interaksi SQL.
* `"fmt"`: Library pemformatan string (untuk `psqlInfo`).
* `"log"`: Library untuk mencatat log error fatal (`log.Fatal`).
* `_ "github.com/lib/pq"`: *Side-effect import* untuk mendaftarkan driver PostgreSQL ke Go.

#### B. Bedah Fungsi (Function) di File Ini:
* **Fungsi: `main()`**
  * **Tujuan**: Fungsi utama yang pertama kali dieksekusi oleh Go saat program dijalankan. Menghubungkan database, melakukan migrasi tabel, dan menyalakan web server.
  * **Bedah Isi Fungsi Baris demi Baris**:
    1. `var PORT = ":50000"`: Menentukan port tempat server HTTP berjalan (`port 50000`).
    2. `psqlInfo := fmt.Sprintf(...)`: Merangkai teks koneksi (*connection string*) PostgreSQL berdasarkan konstanta konfigurasi (`host`, `port`, `user`, `password`, `dbName`).
    3. `DB, err = sql.Open("postgres", psqlInfo)`: Membuka koneksi awal ke database PostgreSQL.
    4. `if err != nil { log.Fatal(...) }`: Jika koneksi gagal dibuat, program langsung dihentikan oleh `log.Fatal`.
    5. `database.DBMigration(DB)`: Memanggil fungsi migrasi tabel otomatis di folder `databases` dengan mengirimkan parameter koneksi `DB`.
    6. `defer DB.Close()`: Menjadwalkan penutupan koneksi database secara otomatis ketika program utama berhenti.
    7. `routers.StartServer().Run(PORT)`: Menyalakan web server HTTP Gin. Perintah `.Run()` bersifat *blocking* (menahan program agar terus aktif mendengarkan request dari pengguna).

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
