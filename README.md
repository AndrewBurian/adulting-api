# Adulting API
API endpoints for recieving and spending adult points

### Routes
/auth
/auth/password  [POST]  login requests
/auth/logout    [POST]  invalidate current API token
/activities     [GET]   get list of activities
/users          [POST]  create new user account
/users/:name    [GET]   information about a user