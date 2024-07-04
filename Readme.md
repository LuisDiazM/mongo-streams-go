# Streams Mongo

## Introducción
En mongo existe un concepto llamado _change streams_ que permite a las aplicaciones acceder a los cambios de la base de datos en tiempo real, este concepto podría compararse con los triggers de firestore donde se pueden escuchar los cambios en una colección de la base de datos.

Este repositorio tiene como objetivo crear una copia de un documento basado en sus cambios de tiempo real como se muestra en el siguiente esquema.

![](/docs/img/stream.jpg)

## Requerimientos
- Una aplicación puede ser un CLI, Service que escuche cambios de escritura de un nuevo documento y actualizaciones en sus campos.
- Crear una copia de los documentos en una colección nueva de mongo.

## Limitaciones
- _Change streams_ está disponible para replicaSet y sharded clusters, en instancias locales de docker no es posible usar esta funcionalidad.
- Este servicio no escala horizontalmente, es decir, si quiero escuchar los cambios de la colección A, sólo puedo tener un stream ya que si se tiene otras instancias del stream escuchando a la misma colección A todas van a recibir los mismos cambios, adicional a ello van a consumir conexiones del pool que tiene mongo donde una mala administración puede ocasionar problemas al cluster.
- Sólo realiza operaciones de insert y update, adicional a ello no aplica cambios dentro de arrays de mongo

## Referencias
- [change streams](https://www.mongodb.com/docs/manual/changeStreams/)