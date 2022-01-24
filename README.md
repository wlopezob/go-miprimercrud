# Mi Primer Crud - GO - Postgres

## Instalar driver Postgres
go get -u github.com/lib/pq

## Modificar los datos de la conexion de BD
```
func conexionBD() *sql.DB {
	driver := "postgres"
	user := "miusuario"
	password := "mipass"
	bd := "mibd"
	server := "miserver"
    ...
```

## Ejecutar
go run src/main.go