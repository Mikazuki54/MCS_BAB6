# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode Sangat Mendalam, Baris demi Baris, Lengkap Tanpa Ada yang Terlewat Termasuk Penjelasan Library / Folder yang Di-import.*

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Import**:
  * *File ini tidak memerlukan library tambahan (`import`), karena murni hanya mendeklarasikan tipe data dasar bentukan Go (*struct*).*
* **Penjelasan Kode**:
  * ◦ `package entities` : Mendeklarasikan bahwa file ini merupakan bagian dari package atau kelompok folder bernama `entities`.
  * ◦ `type Card struct` : Perintah untuk membuat cetakan data baru (*custom data type*) bernama `Card`.
  * ◦ `ID string` : Kolom data di dalam cetakan `Card` bernama `ID` dengan tipe data teks (`string`).
  * ◦ `json: "id"` : Label penanda untuk converter JSON agar field `ID` otomatis dicetak menggunakan huruf kecil `"id"`.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Import**:
  * *File `.sql` tidak menggunakan sistem `import` bahasa Go, melainkan menggunakan perintah DDL (Data Definition Language) SQL murni.*
* **Penjelasan Kode**:
  * ◦ `-- +migrate Up` : Komentar direktif yang dibaca oleh `sql-migrate` untuk memberitahu sistem bahwa skrip ini dipakai untuk membuat atau memperbarui tabel.
  * ◦ `-- +migrate StatementBegin` & `-- +migrate StatementEnd` : Batas awal dan akhir satu blok transaksi perintah SQL.
  * ◦ `CREATE TABLE card (` : Perintah SQL untuk membuat tabel baru bernama `card`.
  * ◦ `id varchar(20)` : Membuat kolom `id` bertipe teks dinamis berukuran maksimal 20 karakter.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Library / Folder yang Di-import**:
  * `"database/sql"`: Library bawaan Go untuk mengelola koneksi dan interaksi dengan database relational.
  * `"embed"`: Fitur bawaan Go untuk menyertakan file fisik (seperti file migrasi `.sql`) langsung ke dalam file program (*binary*).
  * `"fmt"`: Library bawaan Go untuk memformat dan mencetak teks/pesan ke terminal.
  * `migrate "github.com/rubenv/sql-migrate"`: Library pihak ketiga untuk menjalankan migrasi skema database SQL secara otomatis (diberi alias `migrate`).
* **Penjelasan Kode**:
  * ◦ `//go:embed sql_migrations/*.sql` : Perintah compiler Go untuk membungkus file `.sql` ke dalam variabel di bawahnya.
  * ◦ `var dbMigration embed.FS` : Variabel sistem file virtual penampung file `.sql`.
  * ◦ `var DbConnection *sql.DB` : Variabel global pointer database yang bisa diakses oleh file lain.
  * ◦ `func DBMigration(dbParam *sql.DB)` : Fungsi utama migrasi.
    * **`dbParam`**: Parameter kunci koneksi database aktif yang dikirim dari `main.go`.
  * ◦ `n, errs := migrate.Exec(...)` : Menjalankan migrasi.
    * **`n`**: Angka jumlah file migrasi yang sukses dijalankan.
    * **`errs`**: Pesan error jika migrasi gagal.
  * ◦ `if errs != nil { panic(errs) }` : Menghentikan program secara paksa jika migrasi gagal.
  * ◦ `DbConnection = dbParam` : Menyimpan koneksi database ke variabel global.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Library / Folder yang Di-import**:
  * `"bab6/entities"`: Memanggil package `entities` buatan sendiri di dalam proyek ini agar kita bisa menggunakan cetakan struct `entities.Card`.
  * `"database/sql"`: Memanggil pustaka standar Go untuk mengelola tipe data koneksi database SQL (`*sql.DB`).
* **Penjelasan Kode**:
  * ◦ `db *sql.DB` : Parameter koneksi database yang wajib dikirim agar repository tahu ke mana query harus dikirim.
  * ◦ `GetCards` : Mengambil semua data kartu dengan perintah `SELECT`, lalu melakukan perulangan `rows.Next()` dan `rows.Scan(&data.ID)` untuk memindahkan data dari database ke dalam slice Go.
  * ◦ `defer rows.Close()` : Menutup kursor pembacaan baris database secara otomatis di akhir fungsi agar memori tidak bocor.
  * ◦ `InsertCard`, `UpdateCard`, `DeleteCard` : Menjalankan perintah perubahan data (`INSERT`, `UPDATE`, `DELETE`) menggunakan `db.Exec(...)` dengan parameter aman `$1`, `$2`.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Library / Folder yang Di-import**:
  * `database "bab6/databases"`: Meng-import package `databases` buatan sendiri (diberi alias `database`) untuk mengambil variabel global `database.DbConnection`.
  * `"bab6/entities"`: Meng-import struktur data struct `Card`.
  * `"bab6/repositories"`: Meng-import fungsi-fungsi query SQL (*GetCards, InsertCard, dll*).
  * `"net/http"`: Library bawaan Go untuk kode status HTTP (seperti `http.StatusOK` untuk `200`, `http.StatusInternalServerError` untuk `500`, dll).
  * `"github.com/gin-gonic/gin"`: Library framework Gin untuk menangani konteks request HTTP (`gin.Context`) dan helper JSON (`gin.H`).
* **Penjelasan Kode**:
  * ◦ `c *gin.Context` : Objek konteks Gin untuk membaca data request dari klien dan mengirim response.
  * ◦ `database.DbConnection` : Mengambil koneksi database global yang sudah diinisialisasi di `main.go`.
  * ◦ `c.ShouldBindJSON(...)` : Membaca data JSON kiriman pengguna dari body request.
  * ◦ `c.JSON(...)` : Mengirim balasan akhir ke pengguna dalam format JSON beserta status kodenya.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Library / Folder yang Di-import**:
  * `"bab6/controllers"`: Memanggil package `controllers` untuk menghubungkan URL dengan fungsi handler yang dijalankan.
  * `"time"`: Library bawaan Go untuk pengelolaan waktu (digunakan untuk mengatur durasi kedaluwarsa cache CORS `MaxAge`).
  * `"github.com/gin-contrib/cors"`: Middleware resmi dari Gin untuk mengatur kebijakan CORS (*Cross-Origin Resource Sharing*).
  * `"github.com/gin-gonic/gin"`: Framework web Gin untuk membuat instance engine router (`gin.Engine`).
* **Penjelasan Kode**:
  * ◦ `router := gin.Default()` : Membuat instance router utama Gin.
  * ◦ `router.Use(cors.New(...))` : Mengaktifkan aturan keamanan CORS agar website luar / frontend dapat mengakses API ini.
  * ◦ `router.GET`, `router.POST`, `router.PUT`, `router.DELETE` : Mendaftarkan endpoint URL API dan memetakannya ke fungsi controller yang sesuai.

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

#### Penjelasan Import & Baris demi Baris:
* **Penjelasan Library / Folder yang Di-import**:
  * `database "bab6/databases"`: Meng-import package migrasi database buatan sendiri.
  * `"bab6/routers"`: Meng-import package routers untuk menjalankan `StartServer()`.
  * `"database/sql"`: Library bawaan Go untuk mengelola koneksi database SQL.
  * `"fmt"`: Library bawaan Go untuk memformat teks string (digunakan saat merangkai `psqlInfo`).
  * `"log"`: Library bawaan Go untuk mencetak pesan fatal dan menghentikan program jika terjadi error kritis (`log.Fatal`).
  * `_ "github.com/lib/pq"`: **Side-Effect Import** dengan garis bawah `_`. Mengimpor driver PostgreSQL murni agar terdaftar ke package standar `database/sql` tanpa memanggil fungsinya secara langsung.
* **Penjelasan Kode**:
  * ◦ `const (...)` : Menyimpan data konfigurasi tetap (host, port, user, password, dbName).
  * ◦ `sql.Open("postgres", psqlInfo)` : Membuka gerbang koneksi ke PostgreSQL.
  * ◦ `database.DBMigration(DB)` : Memanggil migrasi tabel otomatis.
  * ◦ `defer DB.Close()` : Menutup koneksi database di akhir siklus hidup program.
  * ◦ `routers.StartServer().Run(PORT)` : Menyalakan web server HTTP Gin pada port `:50000`.

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
