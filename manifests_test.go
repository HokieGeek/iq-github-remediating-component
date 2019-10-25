package main

import "testing"

var testPatch = `@@ -75,6 +75,7 @@
],
"dependencies": {
"body-parser": "^1.19.0",
+ "chalk": "^1.0.0",
"check-dependencies": "^1.1.0",
"clarinet": "^0.12.3",
"colors": "^1.3.3",
@@ -85,7 +86,6 @@
"cors": "^2.8.5",
"dottie": "^2.0.1",
"download": "^7.1.0",
- "errorhandler": "^1.5.1",
"express": "^4.17.1",
"express-jwt": "0.1.3",
"express-rate-limit": "^4.0.1",
@@ -112,7 +112,7 @@
"libxmljs2": "^0.21.3",
"marsdb": "^0.6.11",
"morgan": "^1.9.1",
- "moment": "^2.3.0",
+ "moment": "^2.1.0",
"multer": "^1.4.1",
"node-pre-gyp": "^0.13.0",
"notevil": "^1.3.1 ",
`

func TparsePatchAdditions(t *testing.T) {
	t.Skip("TODO")
}
