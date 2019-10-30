package main

import (
	"reflect"
	"testing"
)

var dummyPatches = map[string]string{
	"package.json": `@@ -75,6 +75,7 @@
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
`,
	"requirements.txt": `@@ -19,7 +19,6 @@ flake8==2.3.0
flask-script==2.0.5
flask-session==0.2.1
flask-sslify==0.1.5
-flask==0.10.1
funcsigs==1.0.2 # via mock
futures==3.0.5 # via s3transfer
gevent==1.2.1
@@ -30,7 +29,7 @@ idna==2.6
ipaddress==1.0.22
isodate==0.5.4 # via python3-saml
itsdangerous==0.24 # via flask
-jinja2==2.9.6 # via flask
+jinja2==2.10 # via flask
jmespath==0.9.2 # via boto3, botocore
lru-dict==1.1.6
lxml==4.2.1 # via xmlsec
@@ -40,6 +39,7 @@ mock==2.0.0
ndg-httpsclient==0.4.2
nose-pathmunge==0.1.2
nose==1.3.3
+openpyxl==2.0.5
pbr==2.0.0 # via mock
pep8==1.5.7
pkgconfig==1.3.1 # via xmlsec`,
	"packages.config": `@@ -2,4 +2,5 @@
<packages>
<package id="jQuery" version="3.1.1" targetFramework="net46" />
<package id="NLog" version="4.3.10" targetFramework="net46" />
-</packages>
\ No newline at end of file
+ <package id="LibGit2Sharp" version="0.20.0" targetFramework="net46" />
+</packages>`,
	"pom.xml": `@@ -71,7 +71,7 @@
<dependency>
<groupId>axis</groupId>
<artifactId>axis</artifactId>
- <version>1.2</version>
+ <version>1.2.1</version>
</dependency>
<dependency>
<groupId>axis</groupId>
@@ -91,7 +91,7 @@
<dependency>
<groupId>commons-fileupload</groupId>
<artifactId>commons-fileupload</artifactId>
- <version>1.2.1</version>
+ <version>1.2.2</version>
</dependency>
<dependency>
<groupId>commons-io</groupId>
@@ -101,7 +101,7 @@
<dependency>
<groupId>commons-collections</groupId>
<artifactId>commons-collections</artifactId>
- <version>3.1</version>
+ <version>3.0</version>
</dependency>
<dependency>
<groupId>commons-digester</groupId>
@@ -144,6 +144,11 @@
<artifactId>log4j</artifactId>
<version>1.2.8</version>
</dependency>
+ <dependency>
+ <groupId>org.bouncycastle</groupId>
+ <artifactId>org.bouncycastle</artifactId>
+ <version>1.55</version>
+ </dependency>
<dependency>
<groupId>wsdl4j</groupId>
<artifactId>wsdl4j</artifactId>`,
	"build.gradle": `@@ -11,11 +11,11 @@ repositories {

dependencies {
compile group: 'javax.activation', name: 'activation', version: '1.1'
- compile group: 'axis', name: 'axis', version: '1.2'
+ compile group: 'axis', name: 'axis', version: '1.2.1'
compile group: 'axis', name: 'axis-saaj', version: '1.2'
compile group: 'axis', name: 'axis-jaxrpc', version: '1.2'
compile group: 'axis', name: 'axis-ant', version: '1.2'
- compile group: 'commons-fileupload', name: 'commons-fileupload', version: '1.2.1'
+ compile group: 'commons-fileupload', name: 'commons-fileupload', version: '1.2.2'
compile group: 'commons-io', name: 'commons-io', version: '1.4'
compile group: 'commons-collections', name: 'commons-collections', version: '3.1'
// compile group: 'xml-apis', name: 'xml-apis', version: '1.4.1'`,
	"go.sum": `@@ -287,6 +287,10 @@ github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+
github.com/stretchr/testify v1.2.2 h1:bSDNvY7ZPG5RlJ8otE/7V6gMiyenm9RtJ7IUVIAoJ1w=
github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
+github.com/syncthing/notify v0.0.0-20190709140112-69c7a957d3e2 h1:6tuEEEpg+mxM82E0YingzoXzXXISYR/o/7I9n573LWI=
+github.com/syncthing/notify v0.0.0-20190709140112-69c7a957d3e2/go.mod h1:Sn4ChoS7e4FxjCN1XHPVBT43AgnRLbuaB8pEc1Zcdjg=
+github.com/syncthing/syncthing v0.10.26 h1:OWbhDmaxDWWavG+qFRDpq+F47q14QVHfOHvZyUPe/5s=
+github.com/syncthing/syncthing v0.10.26/go.mod h1:ylURcZ2CSTAidps8DKU5N8KBOiMyfTUOsWoRUS7b11Y=
github.com/syndtr/goleveldb v1.0.0 h1:fBdIW9lB4Iz0n9khmH8w27SJ3QEJ7+IgjPEwGSZiFdE=
github.com/syndtr/goleveldb v1.0.0/go.mod h1:ZVVdQEZoIme9iO1Ch2Jdy24qqXrMMOU6lpPAyBWyWuQ=
github.com/syndtr/goleveldb v1.0.1-0.20190318030020-c3a204f8e965/go.mod h1:9OrXJhf154huy1nPWmuSrkgjPUtUNhA+Zmy+6AESzuA=`,
	"go.mod": `@@ -33,6 +33,7 @@ require (
		 github.com/satori/go.uuid v1.2.0 // indirect
		 github.com/shirou/gopsutil v2.19.9+incompatible // indirect
		 github.com/spf13/cobra v0.0.5 // indirect
	+	github.com/syncthing/syncthing v0.10.26 // indirect
		 github.com/ziutek/mymysql v1.5.4 // indirect
		 go.uber.org/atomic v1.4.0 // indirect
		 go.uber.org/multierr v1.1.0 // indirect
	`,
	"Gemfile": `@@ -44,7 +44,7 @@ gem 'omniauth-saml', '~> 1.10'
 gem 'omniauth', '~> 1.9'
	 
 gem 'discard', '~> 1.1'
-gem 'doorkeeper', '~> 4.2'
+gem 'doorkeeper', '~> 4.3'
 gem 'fast_blank', '~> 1.0'
 gem 'fastimage'
 gem 'goldfinger', '~> 2.1'`,
}

func TestParsePatchAdditions(t *testing.T) {
	type args struct {
		patch string
	}
	tests := []struct {
		name string
		args args
		want map[changeLocation]string
	}{
		/*
			{
				"npm",
				args{dummyPatches["package.json"]},
				map[changeLocation]string{
					changeLocation{Position: 4, Line: 78}:   ` "chalk": "^1.0.0",`,
					changeLocation{Position: 21, Line: 115}: ` "moment": "^2.1.0",`,
				},
			},
		*/
		{
			"nuget",
			args{dummyPatches["packages.config"]},
			map[changeLocation]string{
				changeLocation{Position: 6, Line: 5}: ` <package id="LibGit2Sharp" version="0.20.0" targetFramework="net46" />`,
				changeLocation{Position: 7, Line: 6}: `</packages>`,
			},
		},
		/*
			{
				"pypi",
				args{dummyPatches["requirements.txt"]},
				map[changeLocation]string{
					changeLocation{Position: 13, Line: 32}: `jinja2==2.10 # via flask`,
					changeLocation{Position: 21, Line: 42}: `openpyxl==2.0.5`,
				},
			},
			{
				"gradle",
				args{dummyPatches["build.gradle"]},
				map[changeLocation]string{
					changeLocation{Position: 4, Line: 14}: ` compile group: 'axis', name: 'axis', version: '1.2.1'`,
					changeLocation{Position: 9, Line: 18}: ` compile group: 'commons-fileupload', name: 'commons-fileupload', version: '1.2.2'`,
				},
			},
			{
				"ruby",
				args{dummyPatches["Gemfile"]},
				map[changeLocation]string{
					changeLocation{Position: 5, Line: 47}: `gem 'doorkeeper', '~> 4.3'`,
				},
			},
		*/
		// TODO: golang tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePatchLineAdditions(tt.args.patch); !reflect.DeepEqual(got, tt.want) {
				t.Error("parsePatchLineAdditions()")
				t.Errorf(" Got: %#v\n", got)
				t.Errorf("Want: %#v\n", tt.want)
			}
		})
	}
}

func Test_getMavenComponents(t *testing.T) {
	type args struct {
		patch string
	}
	tests := []struct {
		name    string
		args    args
		want    map[changeLocation]component
		wantErr bool
	}{
		{
			"maven",
			args{dummyPatches["pom.xml"]},
			map[changeLocation]component{
				changeLocation{Position: 5, Line: 74}:   component{format: "maven", group: "axis", name: "axis", version: "1.2.1"},
				changeLocation{Position: 14, Line: 94}:  component{format: "maven", group: "commons-fileupload", name: "commons-fileupload", version: "1.2.2"},
				changeLocation{Position: 23, Line: 104}: component{format: "maven", group: "commons-collections", name: "commons-collections", version: "3.0"},
				changeLocation{Position: 34, Line: 150}: component{format: "maven", group: "org.bouncycastle", name: "org.bouncycastle", version: "1.55"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMavenComponents(tt.args.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMavenComponents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Error("getMavenComponents()")
				t.Errorf(" Got: %v\n", got)
				t.Errorf("Want: %v\n", tt.want)
			}
		})
	}
}
