#Requisites
- go 1.14+
- postgres
####install go:
- Linux https://dl.google.com/go/go1.14.1.windows-amd64.msi 
- Mac OS https://dl.google.com/go/go1.14.1.darwin-amd64.pkg
- Windows https://dl.google.com/go/go1.14.1.darwin-amd64.pkg

###Install and configure Postgresql:
https://www.thegeekstuff.com/2009/04/linux-postgresql-install-and-configure-from-source/

#Install
Clone this repository:``git clone http://wankar.com:3000/kaypi/kaypi_back_geo.git``

Go to the file:
``
cd kaypi_back_geo/
``

Install go dependencies: 
``
go mod vendor
``

#Run Test
Run this command int the application directory (environment vars are necessary):
``
go test ./
``

#Run Application
##Manually
Export the environment vars to connect with the postgresql database:

	 DB_USER //postgres username
	 DB_PASS //user password
	 DB_NAME //name of the database
	 DB_HOST // host
	 DB_SSL_MODE //ssl mode
	 ENDPOINT // endpoint of the application Ex( /location )
	 SERVER_ADDRESS server address Ex( localhost:8080 )
	 
Run the fallow command to run the application:	 ``go run ./``
##Using the run.sh script
Open the `run.sh` script and change the values of the environment vars and run `sh run.sh`.

#GraphQl
##Types
User: 
```
type User{ // User represent an User of the system
      id: ID // User id
      name: String // User name
      category: category // User category
      geo_cord: geo_cord // User geolocation
}
```

geo_cord:
```
type geo_cord { // geo_cord represent a golocation point
    latitude: Float // geolocation latitude
    longitude: Float // geolocation longitude
}
```

category:
```
enum category{ // category represent the category of the user in the system.
    CLIENT // represent a client User
    SERVICE_PROVIDER // represent a Service provider User
    GENERIC // includes both CLIENT users and SERVICE_PORVIDER users (only used to make queries like allUser(category)). 
     // Users cannot be of the GENERIC type,
     // declaring a user of the GENERIC type could return an error.
}
```

##Query
AllUsers: returns all users that match a category.

    allUsers(
        category: Category = GENERIC
      ): [User]
 
 
User: returns an User that match an id.

    user( 
        id: Int 
      ): User
    
 GetCloseUsers: returns the User list close to a position. Use `category=GENERIC` to get 
`CLIENT` and `SERVICE_PROVIDER` users. To more info about the resolution param go to: https://h3geo.org/#/documentation/core-library/resolution-table",

    getCloseUsers( 
        originLat: // latitude of the origin point
        originLong: Float // longitude of the origin point
        resolution: Int  // H3 index resolution. value range (0-15)
        category: category  // category of Users 
      ): [User]

##Mutation
AddUser: set a new User.
    
    addUser(
        name: String
        lat: Float
        long: Float
        category: category
      ): User
      
DeleteUser: remove an User by its Id.

     deleteUser(id: Int): User
     
UpdateGeoCord: update the the geographical coordinates of the User.

    updateGeoCord(
        id: Int
        newLat: Float
        newLong: Float
      ): User 