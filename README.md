# cssmaker
Utility that track directory changes and recompile sass/scss source files to css

## How To Install   
```bash
go install github.com/night-codes/cssmaker
```

## How to use

```bash
cssmaker -i scss/style.scss -o css/style.css -t scss/ 
```

## Options 
```bash  
  -i string
        Input SCSS file. (default "scss/style.scss")
  -o string
        Output CSS file. (default "css/style.css")
  -t string
        Track directory    
```