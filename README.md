# NewsAggregator

## Build

Run `go get` to get dependencies (may need [gcc](http://tdm-gcc.tdragon.net/download) for this)  
Run `go build` to build the project, it will generate an exe file

## Run

*exe_name* -config="*path_to_config_file*" -db="*path_to_database_file*"  
where *path_to_config_file* is path to json file with parsing rules (default: config.json located in the same directory as the exe file),  
      *path_to_database_file* is path to database (default: database.db located in the same directory as the exe file)  

Config schema:  
```json
{  
    "rss": [  
        {  
            "url": "string",  
            "interval": "number"  
        }  
    ],  
    "html": [  
        {  
            "url": "string",  
            "interval": "number",  
            "article_selector": "string",  
            "title_selector": "string",  
            "description_selector": "string"  
        }  
    ]  
}
```