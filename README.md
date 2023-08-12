# simple-article-web-service
A simple article web service with DDD-CQRS 

## Feature
* DDD-CQRS Design
* [Go Fiber](https://docs.gofiber.io/)
* Dockerfile
* Docker Compose
* Postgres
* pgAdmin
* Swagger UI
* Redis
* Grafana, Prometheus, Redis-Exporter
* SEO friendly optimization

## Architecture
![architecture](/doc/architecture.jpg)

## Enpoint
* POST /articles <i>Create a new article</i>
* GET /articles <i>Get list of articles based on title, body and author filters</i>
* GET /articles/{id}/{name} <i>Get an article by ID</i> - {name} for SEO friendly optimization

## Main Folder DDD-CQRS
![main folder](/doc/main-folder-ddd-cqrs.jpg)

## Instalation
```
$ git clone git@github.com:anang5u/simple-article-web-service.git
$ cd simple-article-web-service
$ docker compose up
```
```
PENTING, Sebelum intalasi dilakukan pastikan di localhost telah meng-alokasikan port dibawah ini agar tidak terjadi conflict
* 8999 - Simple article web service API
* 8090 - Swagger UI
* 5050 - pgAdmin
* 3210 - Graphana
* 9090 - Prometheus
* 9121 - Redis Exporter
```
tunggu proses instalasi beberapa saat hingga muncul seperti terlihat pada gambar di bawah ini:

![simple-web-service-ready](/doc/simple-web-service-ready.jpg)

Untuk memulai request POST/GET articles terhadap <b>[simple-article-web-service](http://localhost:8999/articles)</b>, buka Swagger UI pada browser [http://localhost:8090/](http://localhost:8090/) yang akan terlihat seperti pada gambar berikut ini

![Swagger UI](/doc/swagger-ui.jpg)

Silahkan untuk melakukan <i>Create a new article</i> terlebih dahulu sebelum [GET list of article](http://localhost:8999/articles) dan [GET article by ID](http://localhost:8999/articles/1/)</i>. Article yang di ambil dari database atau cache bisa terlihat pada bagian log <i>web service</i>

```
$ docker logs simple-article-web-service-web-1
```
atau untuk melihat semua container
```
$ docker container ls
$ docker container ls -a
```

Contoh Request-Responses Swagger UI untuk enpoint POST /articles <i>Create a new article</i>

![Create a new article](/doc/swagger-ui-create-article.jpg)

Contoh Request-Responses Swagger UI untuk enpoint GET /articles <i>Get list of articles based on title, body and author filters</i>

![Get list articles](/doc/swagger-ui-get-list-articles.jpg)

Contoh Request-Responses Swagger UI untuk enpoint GET /articles/{id}/{name} <i>Get an article by ID</i> - {name} for SEO friendly optimization

![Get article by ID](/doc/swagger-ui-get-article-byid.jpg)

Monitoring tool untuk Postgres, sebagai alternative disini menggunakan <b>[pgAdmin](http://localhost:5050/)</b> yang berada pada url [http://localhost:5050/](http://localhost:5050/), untuk login pgAdmin silahkan menggunakan
```
Email Address/Username: pgadmin@demo.com
Password: password
```

![pgAdmin 4](/doc/pgadmin-4.jpg)

Setelah login, pilih <b>Add New Server</b>
* tab <b>General -> Name</b> : WS DB (bebas)
* tab <b>Connection -> Hostname/address</b> : pgsql-server
* tab <b>Connection -> Username</b> : postgres
* tab <b>Connection -> Password</b> : secr3tPWD

![pgAdmin Create New Server](/doc/pgadmin-4-new-server.jpg)

Kemudian <b>Save</b>, akan muncul tampilan seperti pada gambar berikut

![pgAdmin Monitoring Tools](/doc/pgadmin-monitoring.jpg)

Untuk monitoring Redis sebagai cache, digunakan tool <b>Grafana, Prometheus dan Redis Exporter</b>. Untuk melakukan setting Grafana, pastikan terlebih dahulu bahwa prometheus sudah berjalan sebagai mana mestinya. Lakukan pengecekan pada url [http://localhost:9090](http://localhost:9090/targets?search=), tampilan pada browser untuk prometheus akan terlihat seperti pada gambar berikut ini

![Prometheus](/doc/prometheus.jpg)

Selanjutnya lakukan setting pada Grafana yang berada pada url [localhost:3210](http://localhost:3210)

![Grafana](/doc/grafana.jpg)

Silahkan melakukan login default pada grafana dengan menggunakan

```
Email or Username: admin
Password: admin
```
Tambahkan <b>Data Sources -> Add data source -> Prometheus</b>

![Grafana data source prometheus](/doc/grafana-data-source-prometheus.png)

Pada bagian HTTP Prometheus server URL

```
Prometheus server URL: http://prometheus:9090
```

![Grafana promethes HTTP](/doc/grafana-prometheus.png)

Kemudian pada bagian halaman paling bawah, klik <b>Save & test</b>

![Grafana sava and test promethes](/doc/grafana-save-test-prometheus.jpg)

Selanjutnya untuk menambahkan Dashboard pada Grafana, pilih <b>Menu -> Dashboard</b> kemudian pilih <b>New -> Import</b>

![Grafana new dashboard](/doc/grafana-new-dashboard.jpg)

pada kolom input <i>Import via grafana.com</i> ketikan id 763, lalu klik <b>Load</b>

```
Import via grafana.com: 763
```

selanjutnya pada kolom Promotheus pilih Prometheus, dan klik <b>Import</b>
```
Promotheus : Promotheus (default)
```
![Grafana import dashboard](/doc/grafana-import-redis-dashboard.jpg)

Done, untuk monitoring tool Redis pada Grafana akan tampak seperti ini

![Grafana redis monitoring](/doc/grafana-redis-monitoring.jpg)

## Uninstall aka Bersih-Bersih
```
$ docker compose down -v
$ docker image rm simple-article-web-service_web
$ docker image rm bitnami/redis
$ docker image rm swaggerapi/swagger-ui
$ docker image rm postgres
$ docker image rm oliver006/redis_exporter
$ docker image rm grafana/grafana
$ docker image rm prom/prometheus
$ docker image rm dpage/pgadmin4
```

atau untuk menghapus <b>SEMUA IMAGE</b> yang tidak digunakan bisa menggunakan perintah berikut
```
$ docker image prune -a
WARNING! This will remove all images without at least one container associated to them.
Are you sure you want to continue? [y/N] y
```