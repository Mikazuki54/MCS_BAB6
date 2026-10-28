# Dokumentasi Lengkap Proyek REST API Card Bridge (`bab6`)

Proyek ini adalah aplikasi **RESTful API** berbasis bahasa pemrograman **Go (Golang)** menggunakan framework **Gin Gonic** dan database **PostgreSQL**. Proyek ini menerapkan pola arsitektur berlapis (**Layered Architecture** / MVC-Repository Pattern) untuk mengelola data kartu (`card`) yang mencakup operasi pembacaan data (`GET`), penambahan data (`POST`), dan penghapusan data (`DELETE`), lengkap dengan sistem migrasi database otomatis menggunakan **`sql-migrate`**.

---

## Daftar Isi

1. [Arsitektur & Alur Kerja Aplikasi](#1-arsitektur--alur-kerja-aplikasi)
2. [Struktur Direktori Proyek](#2-struktur-direktori-proyek)
3. [Penjelasan Rinci Masing-Masing File](#3-penjelasan-rinci-masing-masing-file)
   - [3.1. `go.mod` & `go.sum`](#31-gomod--gosum)
   - [3.2. `entities/card_bridge_entity.go`](#32-entitiescard_bridge_entitygo)
   - [3.3. `databases/sql_migrations/1_initiate.sql`](#33-databasessql_migrations1_initiatesql)
   - [3.4. `databases/database.go`](#34-databasesdatabasego)
   - [3.5. `repositories/card_bridge_repo.go`](#35-repositoriescard_bridge_repogo)
   - [3.6. `controllers/card_bridge_controllers.go`](#36-controllerscard_bridge_controllersgo)
   - [3.7. `routers/card_bridge_routers.go`](#37-routerscard_bridge_routersgo)
   - [3.8. `main.go`](#38-maingo)
4. [Dokumentasi Endpoint API](#4-dokumentasi-endpoint-api)
5. [Panduan Menjalankan Proyek](#5-panduan-menjalankan-proyek)

---

## 1. Arsitektur & Alur Kerja Aplikasi

Proyek ini memisahkan tanggung jawab kode ke dalam beberapa lapisan (layer) agar kode rapi, mudah dipelihara (*maintainable*), dan mudah dikembangkan (*scalable*):

```
[ Client / Postman / Browser ]
              │
              ▼ (HTTP Request: GET / POST / DELETE)
┌─────────────────────────────────────────┐
│ 1. Routers (routers/)                   │  -> Menerima request URL & mengarahkan ke Controller
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│ 2. Controllers (controllers/)           │  -> Mengambil parameter request & mengatur format JSON response
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│ 3. Repositories (repositories/)         │  -> Menjalankan query SQL langsung ke Database menggunakan Entities
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│ 4. Database & Entities                  │  -> Menyimpan koneksi DB (`databases/`) & struktur tabel (`entities/`)
│    (databases/ & entities/)             │
└─────────────────────────────────────────┘
              │
              ▼
      [ PostgreSQL DB ]
```

---

## 2. Struktur Direktori Proyek

```text
testbab6/
├── controllers/
│   └── card_bridge_controllers.go      # Handler HTTP untuk memproses request & mengirim response JSON
├── databases/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql              # File skrip SQL untuk membuat tabel 'card' saat migrasi
│   └── database.go                     # Konfigurasi penyimpanan koneksi DB & eksekusi migrasi SQL
├── entities/
│   └── card_bridge_entity.go           # Definisi struct/model data 'Card'
├── repositories/
│   └── card_bridge_repo.go             # Fungsi-fungsi query SQL (SELECT, INSERT, DELETE) ke tabel 'card'
├── routers/
│   └── card_bridge_routers.go          # Definisi endpoint URL API dan pemetaan ke fungsi controller
├── go.mod                              # Definisi nama module Go ('bab6') dan daftar dependensi/library
├── go.sum                              # Checksum kriptografi untuk menjaga integritas versi library
├── main.go                             # Titik masuk utama (entry point) program
└── README.md                           # Dokumentasi lengkap proyek
```

---

## 3. Penjelasan Rinci Masing-Masing File

---

### 3.1. `go.mod` & `go.sum`

**Lokasi File:** `go.mod` dan `go.sum`

#### Fungsi Utama
File `go.mod` adalah file konfigurasi utama modul Go yang mendefinisikan identitas proyek serta daftar pustaka pihak ketiga (*third-party libraries*) yang dibutuhkan aplikasi.

#### Penjelasan Bagian Penting di `go.mod`
1. **`module bab6`**
   Mendeklarasikan bahwa nama module proyek ini adalah `bab6`. Karena nama module-nya adalah `bab6`, maka setiap kali satu folder ingin meng-import package dari folder lain di dalam proyek ini, awalan path import-nya wajib menggunakan `"bab6/<nama_folder>"` (contohnya `"bab6/entities"`, `"bab6/controllers"`, `"bab6/repositories"`, `"bab6/databases"`, `"bab6/routers"`).
2. **`go 1.27.3`**
   Menunjukkan versi toolchain Go yang digunakan pada proyek.
3. **Dependensi Utama (`require`)**:
   - `github.com/gin-gonic/gin`: Framework web HTTP berkinerja tinggi yang digunakan untuk membuat server REST API, membaca parameter URL, dan mengembalikan response berformat JSON.
   - `github.com/lib/pq`: Driver database PostgreSQL murni untuk Go yang bekerja di balik layar bersama package standar `database/sql`.
   - `github.com/rubenv/sql-migrate`: Library untuk menjalankan migrasi skema database secara terstruktur menggunakan file `.sql`.
4. **File `go.sum`**
   Berisi daftar hash kriptografi (checksum) dari setiap versi package yang diunduh. Go menggunakan file ini untuk memastikan bahwa dependensi yang digunakan tidak berubah atau rusak saat dijalankan di komputer lain.

---

### 3.2. `entities/card_bridge_entity.go`

**Lokasi File:** `entities/card_bridge_entity.go`

#### Kode Lengkap
```go
package entities

type Card struct {
	ID string `json:"id"`
}
```

#### Fungsi Utama
File ini berada di lapisan **Entity (Model)**. Tugasnya adalah mendefinisikan struktur data (`struct`) yang merepresentasikan satu baris data pada tabel `card` di database PostgreSQL sekaligus menentukan format objek ketika dikonversi menjadi JSON.

#### Penjelasan Baris demi Baris
- **`package entities`**: Mendefinisikan bahwa file ini merupakan bagian dari package `entities`. Package ini nantinya di-import oleh `repositories` dan `controllers`.
- **`type Card struct { ... }`**: Mendeklarasikan tipe data bentukan bernama `Card` (huruf awal kapital `C` menandakan tipe ini bersifat *exported* atau publik sehingga bisa diakses dari package lain).
- **`ID string \`json:"id"\``**:
  - `ID`: Nama field di dalam kode Go (huruf kapital agar dapat diakses dari luar package dan dapat dibaca oleh library JSON serializer).
  - `string`: Tipe data dari ID kartu, sesuai dengan tipe kolom `varchar(20)` di database PostgreSQL.
  - `` `json:"id"` ``: Merupakan *struct tag*. Ketika struct `Card` diubah menjadi response JSON oleh Gin (`c.JSON(...)`), nama field `ID` akan ditampilkan sebagai `"id"` (huruf kecil) pada output JSON, misalnya:
    ```json
    {
      "id": "CARD001"
    }
    ```

---

### 3.3. `databases/sql_migrations/1_initiate.sql`

**Lokasi File:** `databases/sql_migrations/1_initiate.sql`

#### Kode Lengkap
```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE card (
    id varchar(20)
);

-- +migrate StatementEnd
```

#### Fungsi Utama
File ini berisi instruksi SQL (DDL - *Data Definition Language*) untuk membangun struktur tabel awal di dalam database PostgreSQL secara otomatis ketika aplikasi pertama kali dijalankan.

#### Penjelasan Baris demi Baris
- **`-- +migrate Up`**: Komentar khusus (direktif) yang wajib dikenali oleh library `sql-migrate`. Penanda ini memberi tahu `sql-migrate` bahwa perintah SQL di bawahnya dijalankan ketika proses migrasi bergerak maju (*Up* / membuat perubahan skema).
- **`-- +migrate StatementBegin` & `-- +migrate StatementEnd`**: Penanda batas awal dan akhir dari satu blok perintah SQL utuh agar dieksekusi sebagai satu kesatuan transaksi oleh `sql-migrate`.
- **`CREATE TABLE card ( id varchar(20) );`**: Perintah SQL untuk membuat tabel bernama `card` yang memiliki satu kolom bernama `id` dengan tipe teks `varchar` sepanjang maksimal 20 karakter.
- **Catatan Mekanisme `sql-migrate`**: Ketika file ini berhasil dijalankan, `sql-migrate` secara otomatis membuat tabel internal bernama `gorp_migrations` di PostgreSQL untuk mencatat bahwa file `1_initiate.sql` sudah pernah dieksekusi. Dengan demikian, saat aplikasi dijalankan ulang di kemudian hari, file migrasi yang sama tidak akan dijalankan dua kali (sehingga tidak memunculkan error *table "card" already exists*).

---

### 3.4. `databases/database.go`

**Lokasi File:** `databases/database.go`

#### Kode Lengkap
```go
package database

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS
var DbConnection *sql.DB

func DBMigrate(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)

	if errs != nil {
		panic(errs)
	}

	DbConnection = dbParam

	fmt.Println("Migrations Succes Applied", n, migrations)
}
```

#### Fungsi Utama
File ini bertanggung jawab untuk:
1. Menyematkan (*embed*) seluruh file `.sql` di dalam folder `sql_migrations` langsung ke dalam binary program Go.
2. Menjalankan migrasi database secara otomatis menggunakan `sql-migrate`.
3. Menyimpan objek koneksi database global (`DbConnection`) agar dapat dipakai oleh lapisan `controllers` dan `repositories`.

#### Penjelasan Baris demi Baris
- **`package database`**: Mendeklarasikan nama package sebagai `database` (meskipun nama foldernya adalah `databases`). Oleh karena itu, saat di-import di file lain, ditulis sebagai `database "bab6/databases"`.
- **Import Package**:
  - `"database/sql"`: Package bawaan Go untuk berinteraksi dengan database SQL.
  - `"embed"`: Package bawaan Go (sejak Go 1.16) untuk menyertakan file statis (seperti file `.sql`) ke dalam kompilasi program.
  - `"fmt"`: Package standar untuk mencetak teks ke layar terminal.
  - `migrate "github.com/rubenv/sql-migrate"`: Meng-import library `sql-migrate` dengan alias `migrate`.
- **`//go:embed sql_migrations/*.sql`**:
  Direktif kompilator Go yang sangat penting. Baris ini menginstruksikan Go untuk membaca semua file berekstensi `.sql` di dalam folder `sql_migrations/` dan memasukkannya ke dalam variabel `dbMigrations` tepat di bawahnya. Tanpa baris ini, `dbMigrations` akan kosong dan menyebabkan `panic: open sql_migrations: file does not exist`.
- **`var dbMigrations embed.FS`**: Variabel sistem file virtual (*embedded file system*) yang menampung isi folder `sql_migrations`.
- **`var DbConnection *sql.DB`**: Variabel global bertipe pointer `*sql.DB`. Variabel ini diawali huruf besar (`D`) agar bersifat publik dan bisa dipanggil oleh `controllers` (`database.DbConnection`) saat memanggil fungsi repository.
- **`func DBMigrate(dbParam *sql.DB)`**: Fungsi yang dipanggil dari `main.go` dengan menerima parameter koneksi database aktif (`dbParam`).
  - **`migrations := &migrate.EmbedFileSystemMigrationSource{ ... }`**: Menentukan sumber file migrasi dari `embed.FS` (`dbMigrations`) dengan akar direktori `"sql_migrations"`.
  - **`n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)`**: Menjalankan migrasi dengan dialek `"postgres"` ke arah `migrate.Up` (menerapkan migrasi baru). Nilai kembalian `n` adalah jumlah file migrasi yang baru saja berhasil dijalankan, sedangkan `errs` berisi error jika proses migrasi gagal.
  - **`if errs != nil { panic(errs) }`**: Jika terjadi kesalahan saat migrasi (misalnya database belum dibuat, kredensial salah, atau sintaks SQL keliru), program langsung dihentikan dengan `panic` dan menampilkan pesan error-nya.
  - **`DbConnection = dbParam`**: Menyimpan objek koneksi database yang sudah teruji ke dalam variabel global `DbConnection`.
  - **`fmt.Println("Migrations Succes Applied", n, migrations)`**: Mencetak pesan ke terminal bahwa migrasi berhasil dijalankan beserta jumlah migrasi (`n`) yang diterapkan.

---

### 3.5. `repositories/card_bridge_repo.go`

**Lokasi File:** `repositories/card_bridge_repo.go`

#### Kode Lengkap
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

#### Fungsi Utama
File ini berada pada lapisan **Repository (Data Access Layer)**. Tugasnya adalah menjadi jembatan langsung antara aplikasi Go dan tabel `card` di database PostgreSQL. Tidak ada logika HTTP atau JSON di sini; file ini murni berisi eksekusi query SQL.

#### Penjelasan Rinci Setiap Fungsi

##### 1. Fungsi `GetCards(db *sql.DB) (result []entities.Card, err error)`
Fungsi ini bertugas mengambil seluruh baris data kartu yang ada di tabel `card`.
- **Parameter**: `db *sql.DB` (koneksi aktif ke PostgreSQL).
- **Named Return Values**: `(result []entities.Card, err error)` — Go secara otomatis menginisialisasi variabel `result` (slice dari `entities.Card`) dan `err` (`error`). Karena menggunakan *named return*, perintah `return` tanpa argumen di dalam fungsi akan otomatis mengembalikan nilai terkini dari `result` dan `err`.
- **`sql := "SELECT * FROM card"`**: Menyusun perintah SQL untuk mengambil semua kolom dan baris dari tabel `card`.
- **`rows, err := db.Query(sql)`**: Menjalankan query pembacaan (`SELECT`) yang mengembalikan banyak baris (`*sql.Rows`). Jika terjadi error saat query, fungsi langsung keluar (`return`).
- **`defer rows.Close()`**: Memastikan kursor hasil query (`rows`) selalu ditutup setelah fungsi `GetCards` selesai dieksekusi, guna mencegah kebocoran memori atau koneksi database (*connection leak*).
- **`for rows.Next()`**: Melakukan perulangan untuk membaca setiap baris hasil query satu per satu:
  - `var data entities.Card`: Membuat variabel penampung sementara untuk satu baris kartu.
  - `err = rows.Scan(&data.ID)`: Menyalin nilai kolom `id` dari baris database saat ini ke dalam alamat memori `&data.ID`. Jika proses pembacaan gagal, fungsi langsung mengembalikan error.
  - `result = append(result, data)`: Menambahkan objek `data` yang berhasil dibaca ke dalam slice `result`.
- **`return`**: Mengembalikan daftar kartu (`result`) beserta status `err` (`nil` jika sukses).

##### 2. Fungsi `InsertCard(db *sql.DB, card entities.Card) (err error)`
Fungsi ini bertugas menambahkan satu data kartu baru ke dalam tabel `card`.
- **Parameter**: Koneksi `db *sql.DB` dan objek `card entities.Card` yang berisi `ID` kartu yang ingin disimpan.
- **`sql := "INSERT INTO card(id) values($1)"`**:
  Menyusun query SQL `INSERT`. Simbol `$1` adalah *positional parameter placeholder* khas PostgreSQL untuk parameter pertama. Penggunaan `$1` (prepared statement parameter) sangat penting untuk mencegah serangan **SQL Injection** dibandingkan menggabungkan string secara manual.
- **`_, err = db.Exec(sql, card.ID)`**: Menjalankan perintah SQL yang tidak mengembalikan baris data (`INSERT`). Nilai `card.ID` dimasukkan menggantikan `$1`. Hasil metadata eksekusi (`sql.Result`) diabaikan dengan `_`, dan hanya `err` yang ditangkap.
- **`return err`**: Mengembalikan `nil` jika penyimpanan berhasil, atau pesan error jika gagal.

##### 3. Fungsi `DeleteCard(db *sql.DB, card entities.Card) (err error)`
Fungsi ini bertugas menghapus data kartu dari tabel `card` berdasarkan `id`-nya.
- **Parameter**: Koneksi `db *sql.DB` dan objek `card entities.Card` yang menyimpan `ID` kartu yang akan dihapus.
- **`sql := "DELETE FROM card WHERE id = $1"`**: Query SQL untuk menghapus baris pada tabel `card` yang kolom `id`-nya cocok dengan parameter `$1`.
- **`_, err = db.Exec(sql, card.ID)`**: Mengeksekusi perintah `DELETE` dengan mengirimkan `card.ID` sebagai nilai `$1`.
- **`return err`**: Mengembalikan status error hasil eksekusi.

---

### 3.6. `controllers/card_bridge_controllers.go`

**Lokasi File:** `controllers/card_bridge_controllers.go`

#### Kode Lengkap
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
		c.JSON(http.StatusInternalServerError, gin.H{"errpr": err.Error()})
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
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data is been delete", "id": idCard})

}
```

#### Fungsi Utama
File ini berada pada lapisan **Controller (Handler Layer)**. Tugasnya adalah:
1. Menerima HTTP request yang dikirim oleh pengguna melalui `*gin.Context`.
2. Mengambil parameter dari URL (seperti `:id`).
3. Memanggil fungsi yang sesuai di `repositories` dengan mengirimkan koneksi `database.DbConnection`.
4. Mengemas hasil atau pesan error ke dalam format JSON dan mengirimkannya kembali ke klien beserta kode status HTTP (`200 OK` atau `500 Internal Server Error`).

#### Penjelasan Rinci Setiap Fungsi

##### 1. Fungsi `GetCards(c *gin.Context)`
Menangani request `GET /cards` untuk menampilkan seluruh kartu.
- **`var result gin.H`**: `gin.H` adalah singkatan (alias) dari `map[string]interface{}` yang disediakan oleh Gin untuk mempermudah pembuatan objek JSON dinamis.
- **`card, err := repositories.GetCards(database.DbConnection)`**: Memanggil fungsi `GetCards` milik package `repositories` menggunakan koneksi database global `database.DbConnection`.
- **Percabangan `if err != nil`**:
  - Jika terjadi error, `result` diisi dengan pesan error: `{"result": "<pesan error>"}`.
  - Jika berhasil, `result` diisi dengan slice data kartu: `{"result": [...]}`.
- **`c.JSON(http.StatusOK, result)`**: Mengirimkan response HTTP berstatus `200 OK` (`http.StatusOK`) dengan isi body berupa objek `result` yang telah dikonversi menjadi JSON.

##### 2. Fungsi `InsertCard(c *gin.Context)`
Menangani request `POST /card/input/:id` untuk menambahkan kartu baru.
- **`var card entities.Card`**: Membuat instance kosong dari struct `Card`.
- **`idCard := c.Param("id")`**: Mengambil nilai parameter dinamis `:id` dari URL. Contohnya, jika klien mengakses `POST /card/input/KARTU01`, maka `idCard` akan bernilai `"KARTU01"`.
- **`card.ID = idCard`**: Memasukkan nilai `idCard` ke dalam field `ID` pada struct `card`.
- **`err := repositories.InsertCard(database.DbConnection, card)`**: Memanggil repository untuk menyimpan data kartu ke PostgreSQL.
- **`if err != nil`**: Jika terjadi kegagalan penyimpanan ke database, controller mengirimkan response HTTP `500 Internal Server Error` beserta detail error-nya.
- **`c.JSON(http.StatusOK, card)`**: Mengirimkan response HTTP `200 OK` dengan menampilkan kembali objek `card` yang baru ditambahkan dalam format JSON (`{"id": "KARTU01"}`).

##### 3. Fungsi `DeleteCard(c *gin.Context)`
Menangani request `DELETE /card/delete/:id` untuk menghapus kartu berdasarkan ID.
- **`var card entities.Card`**: Membuat variabel bertipe `entities.Card`.
- **`idCard := c.Param("id")`**: Membaca parameter `:id` dari path URL (misalnya `/card/delete/KARTU01` menghasilkan `idCard = "KARTU01"`).
- **`card.ID = idCard`**: Mengisi field `card.ID` dengan `idCard`.
- **`err := repositories.DeleteCard(database.DbConnection, card)`**: Memanggil fungsi `DeleteCard` di repository untuk menghapus baris dengan ID tersebut dari database.
- **`if err != nil`**: Jika query hapus gagal, mengirimkan status `500 Internal Server Error` beserta alasan error-nya.
- **`c.JSON(http.StatusOK, gin.H{"message": "Data is been delete", "id": idCard})`**: Mengirimkan status `200 OK` dengan pesan konfirmasi bahwa data dengan ID tersebut telah dihapus.

---

### 3.7. `routers/card_bridge_routers.go`

**Lokasi File:** `routers/card_bridge_routers.go`

#### Kode Lengkap
```go
package routers

import (
	"bab6/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### Fungsi Utama
File ini berada pada lapisan **Router**. Tugasnya adalah menginisialisasi engine HTTP Gin dan mendaftarkan daftar alamat URL (*endpoint*) beserta metode HTTP-nya (`GET`, `POST`, `DELETE`) agar diarahkan ke fungsi handler yang tepat di dalam `controllers`.

#### Penjelasan Baris demi Baris
- **`package routers`**: Mendefinisikan package `routers`.
- **`import "bab6/controllers"` & `"github.com/gin-gonic/gin"`**: Meng-import package `controllers` dari proyek sendiri dan framework `gin`.
- **`func StartServer() *gin.Engine`**: Fungsi publik yang mengembalikan pointer `*gin.Engine` (instance server router Gin yang siap dijalankan).
- **`router := gin.Default()`**: Membuat instance router Gin dengan dua *middleware* bawaan yang langsung aktif secara otomatis:
  1. **Logger Middleware**: Mencatat setiap request HTTP yang masuk ke terminal (metode, path, status code, dan waktu eksekusi).
  2. **Recovery Middleware**: Menangkap *panic* jika terjadi kesalahan tak terduga saat melayani request agar server tidak langsung mati (*crash*), melainkan mengembalikan status `500 Internal Server Error`.
- **Pendaftaran Route Endpoint**:
  - **`router.GET("/cards", controllers.GetCards)`**: Mendaftarkan endpoint metode `GET` pada path `/cards` untuk dijalankan oleh fungsi `controllers.GetCards`.
  - **`router.POST("/card/input/:id", controllers.InsertCard)`**: Mendaftarkan endpoint metode `POST` pada path `/card/input/:id` (di mana `:id` adalah *path parameter* dinamis) untuk dijalankan oleh fungsi `controllers.InsertCard`.
  - **`router.DELETE("/card/delete/:id", controllers.DeleteCard)`**: Mendaftarkan endpoint metode `DELETE` pada path `/card/delete/:id` untuk dijalankan oleh fungsi `controllers.DeleteCard`.
- **`return router`**: Mengembalikan objek `router` yang sudah dikonfigurasi ke pemanggilnya di `main.go`.

---

### 3.8. `main.go`

**Lokasi File:** `main.go`

#### Kode Lengkap
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
	dbName   = "mcs_bab_6"
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

	database.DBMigrate(DB)

	defer DB.Close()

	routers.StartServer().Run(PORT)
	fmt.Printf("Connect Successfully")
}
```

#### Fungsi Utama
`main.go` adalah **titik masuk utama (Entry Point)** dari seluruh aplikasi. Saat perintah `go run main.go` dijalankan, fungsi `main()` di dalam file inilah yang dieksekusi pertama kali oleh Go untuk menyambungkan database, menjalankan migrasi tabel, dan menyalakan web server HTTP.

#### Penjelasan Baris demi Baris
- **`package main`**: Di Go, program yang dapat dieksekusi (*executable*) wajib berada di dalam `package main` dan memiliki fungsi `main()`.
- **Import Package**:
  - `database "bab6/databases"`: Meng-import package `database` untuk memanggil fungsi migrasi `DBMigrate`.
  - `"bab6/routers"`: Meng-import package `routers` untuk memanggil `StartServer()`.
  - `"database/sql"`, `"fmt"`, `"log"`: Package standar Go untuk koneksi SQL, pemformatan string, dan pencatatan log error.
  - **`_ "github.com/lib/pq"`**: Menggunakan *blank identifier* (`_`). Ini disebut **side-effect import**. Kita tidak memanggil fungsi dari package `pq` secara langsung di dalam kode `main.go`, tetapi kita butuh fungsi `init()` di dalam `github.com/lib/pq` berjalan agar mendaftarkan driver bernama `"postgres"` ke dalam package standar `database/sql`. Tanpa baris ini, `sql.Open("postgres", ...)` akan gagal dengan error *unknown driver "postgres"*.
- **Blok `const (...)`**:
  Mendefinisikan konstanta konfigurasi koneksi ke server PostgreSQL lokal:
  - `host = "localhost"`: Alamat server database (komputer lokal).
  - `port = 5432`: Port default PostgreSQL.
  - `user = "postgres"`: Username database PostgreSQL.
  - `password = "nickytirta54"`: Password user PostgreSQL.
  - `dbName = "mcs_bab_6"`: Nama database yang dituju di PostgreSQL.
- **Blok `var ( DB *sql.DB; err error )`**: Mendeklarasikan variabel global untuk menampung objek koneksi database (`DB`) dan objek `err`.
- **Di dalam `func main()`**:
  1. **`var PORT = ":50000"`**: Menentukan port tempat server HTTP Gin akan berjalan, yaitu port `50000` (`http://localhost:50000`).
  2. **`psqlInfo := fmt.Sprintf(...)`**: Merangkai *Connection String* (DSN - *Data Source Name*) PostgreSQL menjadi:
     `host=localhost port=5432 user=postgres password=nickytirta54 dbname=mcs_bab_6 sslmode=disable`
     (`sslmode=disable` digunakan karena koneksi dilakukan secara lokal tanpa sertifikat SSL).
  3. **`DB, err = sql.Open("postgres", psqlInfo)`**: Menginisialisasi objek koneksi database `*sql.DB` menggunakan driver `"postgres"`. Jika format konfigurasi salah, `log.Fatal` akan mencetak error dan menghentikan program.
  4. **`database.DBMigrate(DB)`**: Memanggil fungsi `DBMigrate` dari `databases/database.go`. Pada tahap ini koneksi ke PostgreSQL benar-benar diuji, tabel `card` dibuat (jika belum ada) melalui `1_initiate.sql`, dan koneksi disimpan ke variabel `database.DbConnection`.
  5. **`defer DB.Close()`**: Menjadwalkan penutupan koneksi database ketika fungsi `main()` selesai atau berhenti.
  6. **`routers.StartServer().Run(PORT)`**: Memanggil `StartServer()` dari `routers/card_bridge_routers.go` lalu menjalankan HTTP server pada port `:50000`. Fungsi `.Run(PORT)` bersifat *blocking* (terus berjalan mendengarkan request HTTP sampai program dihentikan dengan `Ctrl+C`).

---

## 4. Dokumentasi Endpoint API

Server berjalan secara default di alamat: **`http://localhost:50000`**

### 1. Menambahkan Kartu Baru (`POST`)
- **URL**: `/card/input/:id`
- **Method**: `POST`
- **Contoh Request (cURL)**:
  ```bash
  curl -X POST http://localhost:50000/card/input/CARD-001
  ```
- **Contoh Response Sukses (`200 OK`)**:
  ```json
  {
    "id": "CARD-001"
  }
  ```

### 2. Mengambil Semua Data Kartu (`GET`)
- **URL**: `/cards`
- **Method**: `GET`
- **Contoh Request (cURL)**:
  ```bash
  curl -X GET http://localhost:50000/cards
  ```
- **Contoh Response Sukses (`200 OK`)**:
  ```json
  {
    "result": [
      {
        "id": "CARD-001"
      }
    ]
  }
  ```

### 3. Menghapus Data Kartu Berdasarkan ID (`DELETE`)
- **URL**: `/card/delete/:id`
- **Method**: `DELETE`
- **Contoh Request (cURL)**:
  ```bash
  curl -X DELETE http://localhost:50000/card/delete/CARD-001
  ```
- **Contoh Response Sukses (`200 OK`)**:
  ```json
  {
    "id": "CARD-001",
    "message": "Data is been delete"
  }
  ```

---

## 5. Panduan Menjalankan Proyek

1. **Pastikan PostgreSQL Aktif dan Database Sudah Dibuat**
   Buat database bernama `mcs_bab_6` di PostgreSQL (melalui pgAdmin atau `psql`):
   ```sql
   CREATE DATABASE mcs_bab_6;
   ```
2. **Sesuaikan Kredensial Database (Jika Perlu)**
   Periksa konstanta `host`, `port`, `user`, `password`, dan `dbName` di dalam file `main.go` agar sesuai dengan pengaturan PostgreSQL di komputer Anda.
3. **Unduh dan Rapikan Dependensi Go**
   ```bash
   go mod tidy
   ```
4. **Jalankan Server Aplikasi**
   ```bash
   go run main.go
   ```
5. **Tampilan Terminal Saat Berhasil Berjalan**
   ```text
   Migrations Succes Applied 1 &{...}
   [GIN-debug] GET    /cards                    --> bab6/controllers.GetCards (3 handlers)
   [GIN-debug] POST   /card/input/:id           --> bab6/controllers.InsertCard (3 handlers)
   [GIN-debug] DELETE /card/delete/:id          --> bab6/controllers.DeleteCard (3 handlers)
   [GIN-debug] Listening and serving HTTP on :50000
   ```

