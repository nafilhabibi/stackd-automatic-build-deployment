# 🚀 stackd

**stackd** adalah sebuah *Command-Line Interface (CLI)* berbasis bahasa Go yang diciptakan untuk menyederhanakan proses *deployment* project web Anda di lingkungan server Linux.

Hanya dengan menyebutkan path direktori *source code* project Anda, **stackd** akan beraksi layaknya *CI/CD pipeline*: mendeteksi otomatis *environment* yang digunakan—seperti **Node.js, Python, Laravel (PHP), atau Golang**—, melakukan instalasi dependensi atau *virtual environment* yang dibutuhkan, meracik perintah eksekusi yang optimal untuk lingkungan *production*, lalu mengikatnya sepenuhnya ke `systemd` supaya aplikasi Anda bisa langsung berjalan di *background*, termonitor, dan otomatis *restart* setiap kali server dihidupkan ulang.

---

## ✨ Fitur Utama

- 🔍 **Auto-Detection Pintar:** Dapat mendeteksi secara otomatis jenis project dari isi direktorinya (seperti melihat eksistensi `package.json`, `requirements.txt`, `composer.json`, atau `go.mod`).
- 🤖 **Smart CI/CD Initialization:**
  - **Python**: Menganalisa dan mengkonfigurasi *Virtual Environment* (`venv`) secara otomatis serta meng-install pustaka dari `requirements.txt`. Mendukung usulan perintah eksekusi server *production* tangguh secara mandiri (misal: *Gunicorn / Uvicorn*).
  - **Node.js**: Memeriksa folder `node_modules` dan menawarkan automasi proses `npm install` apabila sedang *fresh install*.
  - **Laravel (PHP)**: Memeriksa folder `vendor` dan menawarkan *build* optimal dengan mengeksekusi otomatis `composer install --no-dev --optimize-autoloader`.
- ⚙️ **One-Stop Systemd Builder:** *stackd* membuatkan file `.service` di direktori `/etc/systemd/system/`, melakukan *daemon-reload*, mengatur ekstensi *startup (enable)*, dan menyalakan servis (start) dalam satu ketukan terminal!

---

## 📋 Prasyarat

Sebelum Anda mulai menggunakan **stackd**, usahakan server/device Anda telah memenuhi requirement minimum ini:
- Sistem Operasi **Linux** (terutama tipe *systemd*).
- Memiliki akses *superuser* (perintah `sudo`).
- Tersedia *runtime* environment bahasa terkait (seperti instalasi paket `python3-venv`, `Node/npm`, atau `Composer/PHP` tergantung project jenis yang dikembangkan).

---

## 📥 Instalasi

1. Pastikan Golang versi 1.20+ sudah terinstal di server Anda (jika akan *build* dari source).
2. Clone repository ini (atau unduh source codenya):
   ```bash
   git clone https://github.com/nafilhabibi/stackd-automatic-build-deployment
   cd stackd
   ```
3. Lakukan proses Eksekusi *Installer*:
   Saya telah menyediakan *script shell installer* agar proses kompilasi hingga instalasi global selesai dalam satu ketikan. Anda hanya cukup mengetikkan *command* berikut (dan Anda mungkin akan dimintai *password sudo* server Anda):
   ```bash
   ./install.sh
   ```

   **_Catatan:_** Bila karena suatu hal terminal Anda *permission denied* saat menjalankan skrip, ketikkan dulu *command*: `chmod +x install.sh`
---

## 🛠 Panduan Pemakaian (*Usage*)

Tampilan antarmuka **stackd** sangat interaktif, membimbing pengguna untuk menyepakati atau memperbaiki default secara bertahap.

### Men-Deploy Project Baru:
Langkah pertama yang Anda butuhkan adalah mengeksekusi satu *command* berikut:
```bash
sudo stackd deploy /path/ke/direktori/project-anda
```

### Melihat Daftar Service Anda:
Anda dapat mencetak daftar aplikasi / service apa saja yang pernah Anda deploy di *systemd* server ini secara rapi:
```bash
stackd list
```

### Menghapus (Undeploy) Secara Permanen:
Untuk menghapus service yang telah berjalan secara permanen:
```bash
sudo stackd remove nama-service
```

*(Catatan: pastikan Anda menggunakan `sudo` saat melakukan deploy/remove karena `stackd` butuh permisi.*

### Contoh Penggunaan pada Aplikasi Node.js "Kosong" (Belum `npm install`):

```text
=====================================
🚀 stackd CLI - Deployment Tools
=====================================
🔍 Menganalisa project...
✅ Terdeteksi sebagai  : NODE
...
🎉 Deployment systemd service berhasil!
```

### Troubleshoot (Panduan Jika Masalah Terjadi)

- `Gagal menulis file unit systemd ke /etc/systemd...`: Terjadi jika Anda memanggil `stackd` _tanpa hak permision root / `sudo`_.
- `Gagal membuat venv`: Ini akibat environment instalasi *python / distro* belum menginstall paket venv dasar. Pada keluarga debian/ubuntu misalnya, ketikkan dulu `sudo apt install python3-venv`.

---

## 🤝 Kontribusi

Tool ini adalah *Open-Source* dan siapa saja dapat mengirim *Pull Request* baru! Bila Anda ingin memodifikasi script detektor, menambahkan bahasa framework seperti *Ruby on Rails* atau merancang perbaikan penanganan systemd, silakan fork repositorinya!

## 📜 Lisensi
Aplikasi ini dirilis di bawah naungan **MIT License**. Lihat file [LICENSE](LICENSE) untuk selengkapnya.

> Diciptakan dengan ❤️ untuk menyembuhkan sakit kepala setup server Linux.
