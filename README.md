# Muninn
Simple asset browser web-based frontend.

### About

Simple read-only visualization for the [wysteria](https://github.com/voidshard/wysteria) 
asset tracking & versioning system. Although written with wysteria in mind internally the code is 
kept very generic and interfaces are used for the actual fetching & caching of data from external 
systems, so this can be used to display any similarly-minded system if you had a mind to.

### Building

The simplest way to build the server would be to use 
```bash
go install https://github.com/voidshard/muninn
```

But assuming you have the Go dependencies already kicking about, you could also use
```bash
git clone https://github.com/voidshard/muninn
cd ./muninn
go build -o muninn *.go
```

The client side web page is written in [ReactJS](https://facebook.github.io/react/) and can be built with [npm](https://www.npmjs.com/) 
```bash
cd ./ui
npm run-script build
```
That is, assuming npm in installed.

### Running

Simply
```bash
./muninn
```

By default, this tries to use {current folder}/ui/dist/ as the web root and port 7600 to listen on - but the server
accepts cli args to alter this behaviour if required. 

### Notes

- This doesn't supply any args to the wysteria client, so you may want to check your WYSTERIA_CLIENT_INI env var.
- I include the npm pre-built bundle.js file, but not the js dependencies (mostly cause there's quite a few..). 
- The web client will search for a collection on start so it can display something :) 
- The server caches searches in os.TempDir and the client paginates & loads and unloads table data as you scroll .. so
 it should be able to browse reasonably large datasets. I haven't tested on sets of more than a few thousand entries, so
 I'd be curious to hear about what happens. 

## ToDo
- Add autocomplete magic
