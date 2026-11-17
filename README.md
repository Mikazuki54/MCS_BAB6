# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode Sangat Mendalam, Baris demi Baris, Lengkap Tanpa Ada yang Terlewat.*

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
| **Struct** | Fitur Go untuk membuat wadah tipe data sendiri (seperti objek). |
| **Pointer (`*`)** | Tanda penunjuk alamat memori komputer asli, bukan menduplikat datanya. |
| **Parameter** | Nilai/bahan masukan yang dikirim ke dalam sebuah fungsi agar bisa diolah. |
| **JSON** | Format teks pertukaran data di internet (`"key": "value"`). |
| **CORS** | Sistem keamanan browser agar web lain tidak sembarangan mengambil data API kita. |
| **`defer`** | Perintah "nanti saja" untuk menjalankan fungsi tepat sebelum fungsi utama selesai. |

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

## 📘 5. Penjelasan Rinci Kode Berdasarkan Urutan Alur (`Entity` ➔ `Migration` ➔ `Database` ➔ `Repository` ➔ `Controller` ➔ `Router` ➔ `Main`)

---

### 5.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

**Tugas File:** Mendefinisikan bentuk data kartu yang dipakai dalam program Go dan format teks JSON saat dikirim ke pemanggil API.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Penjelasan Detail Baris demi Baris:
* ◦ `package entities` : Mendeklarasikan bahwa file ini merupakan bagian dari package atau kelompok folder bernama `entities`. Di Go, semua file dalam satu folder wajib memiliki nama package yang sama.
* ◦ `type Card struct` : Perintah untuk membuat cetakan data baru (*custom data type*) bernama `Card`. Huruf `C` besar di awal kata menandakan bahwa struct ini bersifat *exported* (publik) sehingga bisa dibaca dan diakses oleh file di folder lain.
* ◦ `ID string` : Properti atau kolom di dalam cetakan `Card` bernama `ID` dengan tipe data teks (`string`). Digunakan untuk menampung nilai unik kartu.
* ◦ `json: "id"` : *Struct tag* atau label penanda untuk converter JSON bawaan Go. Saat objek `Card` dikonversi menjadi response JSON, field `ID` akan otomatis dicetak menggunakan huruf kecil `"id"`.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `-- +migrate Up` : Komentar khusus (*direktif*) yang dibaca oleh library `sql-migrate`. Perintah ini memberi tahu sistem bahwa kode di bawahnya harus dieksekusi untuk *membuat atau memperbarui* struktur tabel di database.
* ◦ `-- +migrate StatementBegin` & `-- +migrate StatementEnd` : Penanda batas awal dan akhir dari satu blok transaksi perintah SQL utuh agar dieksekusi secara bersamaan oleh migrator.
* ◦ `CREATE TABLE card (` : Perintah SQL dasar untuk membuat tabel baru di dalam database PostgreSQL dengan nama tabel `card`.
* ◦ `id varchar(20)` : Perintah SQL untuk membuat satu kolom di dalam tabel `card` bernama `id` dengan tipe data teks dinamis (`varchar`) yang memiliki panjang maksimal 20 karakter.

---

### 5.3. `databases/database.go` (Pengelola Koneksi DB & Migrasi Otomatis)

**Tugas File:** Membaca file `.sql`, mengeksekusi pembuatan tabel ke PostgreSQL secara otomatis, dan menyimpan koneksi database agar siap dipakai oleh fungsi lain.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `package database` : Menandakan file ini milik kelompok folder `database`.
* ◦ `import (...)` : Mengambil pustaka standar Go (`database/sql`, `embed`, `fmt`) dan pustaka pihak ketiga `sql-migrate` (yang diberi alias `migrate`).
* ◦ `//go:embed sql_migrations/*.sql` : Komentar compiler khusus di Go (Go Embed). Perintah ini menyuruh Go untuk menyertakan file-file `.sql` ke dalam hasil akhir file program (*binary*), sehingga kita tidak perlu membawa folder `sql_migrations` terpisah saat aplikasi dijalankan.
* ◦ `var dbMigration embed.FS` : Variabel bertipe sistem file virtual (*embedded file system*) yang menampung seluruh file migrasi `.sql`.
* ◦ `var DbConnection *sql.DB` : Variabel global bertipe pointer database (`*sql.DB`). Menggunakan huruf awal kapital (`D`) agar bisa diakses dari luar package (misalnya oleh `controllers` atau `repositories`).
* ◦ `func DBMigration(dbParam *sql.DB)` : Fungsi utama untuk menjalankan migrasi database.
  * **Penjelasan `dbParam`**: Parameter/bahan masukan berupa objek koneksi database aktif (`*sql.DB`) yang dikirim langsung dari file `main.go`. Variabel ini berfungsi sebagai jembatan agar fungsi migrasi tahu ke database PostgreSQL mana perintah SQL harus dijalankan.
* ◦ `migrations := &migrate.EmbedFileSystemMigrationSource{...}` : Menyiapkan konfigurasi sumber file migrasi dengan menunjuk ke sistem file virtual `dbMigration` di dalam folder utama `"sql_migrations"`.
* ◦ `n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)` : Perintah utama untuk mengeksekusi file migrasi ke database.
  * **`n`**: Menampung nilai angka (`integer`) yang menunjukkan berapa banyak file migrasi baru yang berhasil dijalankan.
  * **`errs`**: Menampung nilai `error` jika proses migrasi menemui kegagalan.
  * **`dbParam`**: Menyertakan koneksi database yang diterima dari parameter fungsi di atas.
  * **`"postgres"`**: Memberitahu library migrasi bahwa jenis database yang kita tuju adalah PostgreSQL.
  * **`migrate.Up`**: Menandakan arah migrasi maju (membuat tabel baru).
* ◦ `if errs != nil { panic(errs) }` : Pengaman jika terjadi error saat migrasi (misalnya database belum dibuat atau password salah), program akan dihentikan secara paksa (`panic`) agar kita bisa segera memperbaiki kesalahannya.
* ◦ `DbConnection = dbParam` : Menyimpan koneksi database yang sudah teruji dan berhasil dimigrasi ke variabel global `DbConnection`.
* ◦ `fmt.Println(...)` : Mencetak pesan sukses migrasi ke layar terminal komputer.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `package repositories` : Menandakan file ini berada di kelompok khusus repositori data.
* ◦ `db *sql.DB` (pada setiap fungsi repo) : Parameter kunci koneksi database yang dikirim dari controller agar repository bisa berbicara dengan PostgreSQL.
* ◦ **Fungsi `GetCards(db *sql.DB)`**: Mengambil seluruh data kartu.
  * `db.Query(sql)` : Menjalankan perintah `SELECT` untuk mengambil banyak baris data dari database. Hasilnya disimpan di variabel `rows`.
  * `defer rows.Close()` : Menjadwalkan penutupan pointer `rows` di akhir fungsi agar tidak terjadi kebocoran memori RAM (*connection leak*).
  * `for rows.Next()` : Perulangan untuk membaca hasil query baris demi baris.
  * `rows.Scan(&data.ID)` : Perintah untuk menyalin isi kolom database baris saat ini ke dalam variabel memori `&data.ID`.
  * `result = append(result, data)` : Menambahkan data kartu yang baru dibaca ke dalam daftar slice `result`.
* ◦ **Fungsi `InsertCard(db *sql.DB, card entities.Card)`**: Menambahkan kartu baru.
  * `db.Exec(sql, card.ID)` : Menjalankan query yang tidak mengembalikan baris data (`INSERT`/`UPDATE`/`DELETE`). Simbol `$1` digunakan sebagai parameter aman untuk menghindari serangan *SQL Injection*.
* ◦ **Fungsi `UpdateCard(db *sql.DB, card entities.Card)`**: Memperbarui data kartu.
  * Menggunakan perintah SQL `UPDATE card SET id = $1 WHERE id = $2` untuk mengganti nilai ID lama dengan yang baru.
* ◦ **Fungsi `DeleteCard(db *sql.DB, card entities.Card)`**: Menghapus data kartu berdasarkan ID.
  * Menggunakan perintah SQL `DELETE FROM card WHERE id = $1`.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `package controllers` : Menandakan file ini bertugas sebagai pengatur alur kontrol request.
* ◦ `c *gin.Context` : Objek konteks Gin yang membawa seluruh informasi HTTP Request dari klien dan menyediakan alat untuk mengirim response balik.
* ◦ **Fungsi `GetCards(c *gin.Context)`**: Menangani request `GET /cards`.
  * Memanggil `repositories.GetCards(database.DbConnection)` untuk mengambil data dari database.
  * Menggunakan `c.JSON(http.StatusOK, result)` untuk mengirimkan response berupa JSON dengan status HTTP `200 OK`.
* ◦ **Fungsi `InsertCard(c *gin.Context)`**: Menangani request `POST`.
  * `c.Param("id")` : Mengambil nilai parameter `:id` dari URL.
  * `c.JSON(http.StatusInternalServerError, ...)` : Mengirim status `500` jika penyimpanan database gagal.
* ◦ **Fungsi `UpdateCard(c *gin.Context)`**: Menangani request `PUT`.
  * `c.ShouldBindJSON(&card)` : Membaca data JSON yang dikirim pengguna di body request dan menyalinnya ke struct `card`.
  * `http.StatusBadRequest` : Mengembalikan status `400` jika format JSON yang dikirim pengguna tidak valid.
* ◦ **Fungsi `DeleteCard(c *gin.Context)`**: Menangani request `DELETE`.
  * Mengambil ID dari URL, menghapusnya via repository, lalu mengirim pesan sukses `"Data Berhasil Dihapus"`.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `package routers` : Kelompok file rute alamat URL.
* ◦ `router := gin.Default()` : Membuat instance router Gin baru yang sudah dilengkapi dengan logger terminal dan sistem pemulihan error (*recovery*).
* ◦ `router.Use(cors.New(...))` : Konfigurasi CORS (Cross-Origin Resource Sharing).
  * `AllowOrigins: []string{"*"}` : Mengizinkan akses dari domain/frontend mana saja.
  * `AllowMethods` : Menentukan metode HTTP yang diizinkan (`GET`, `POST`, `PUT`, `DELETE`, dll.).
* ◦ `router.GET(...)`, `router.POST(...)`, `router.PUT(...)`, `router.DELETE(...)` : Menghubungkan alamat URL endpoint dengan fungsi controller yang bersangkutan.
* ◦ `return router` : Mengembalikan mesin router yang siap dijalankan ke file `main.go`.

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

#### Penjelasan Detail Baris demi Baris:
* ◦ `package main` : Package utama aplikasi Go yang dapat dijalankan secara mandiri.
* ◦ `_ "github.com/lib/pq"` : *Side-effect import* dengan garis bawah `_`. Berfungsi mendaftarkan driver database PostgreSQL ke dalam package standar `database/sql` Go agar komputer mengenali perintah database.
* ◦ `const (...)` : Kumpulan data konstanta (tetap) untuk alamat host, port, username, password, dan nama database PostgreSQL.
* ◦ `sql.Open("postgres", psqlInfo)` : Perintah untuk membuka gerbang koneksi awal ke database PostgreSQL menggunakan string konfigurasi `psqlInfo`.
* ◦ `database.DBMigration(DB)` : Memanggil fungsi migrasi otomatis yang ada di folder `databases` dengan mengirimkan parameter koneksi `DB`.
* ◦ `defer DB.Close()` : Perintah penjadwalan agar koneksi database ditutup secara otomatis tepat sebelum program dimatikan.
* ◦ `routers.StartServer().Run(PORT)` : Perintah utama untuk menyalakan server HTTP Gin di port `:50000`. Perintah ini bersifat *blocking* (menahan program agar terus berjalan mendengarkan request).

---

## 🚀 8. Cara Menjalankan Proyek Utama (`bab6`)

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

## 🧪 9. Pengujian Endpoint API (Termasuk PUT)

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
