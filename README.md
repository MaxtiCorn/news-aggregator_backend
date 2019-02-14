# NewsAggregator

## Build

Run `go get` to get dependencies (may need [gcc](http://tdm-gcc.tdragon.net/download) for this)  
Run `go build` to build the project, it will generate an exe file

## Run

*exe_name* -config="*path_to_config_file*" -db="*path_to_database_file*" -port="*port*"
where  
*path_to_config_file* is path to json file with parsing rules (default: config.json located in the same directory as the exe file),  
*path_to_database_file* is path to database (default: database.db located in the same directory as the exe file)  
*port* is port to run app (default: 69)  

Config schema:  
```json
{  
    "rss": [  
        {  
            "source": "string",
            "url": "string",  
            "interval": "number"  
        }  
    ],  
    "html": [  
        {  
            "source": "string",
            "host": "string",  
            "url": "string",  
            "interval": "number",  
            "article_selector": "string",  
            "title_selector": "string",  
            "description_selector": "string",  
            "link_selector": "string"
        }  
    ]  
}
```

where *host* is root of site  
*url* is url of page with feed  
*interval* is polling interval  
*article_selector* is css selector to find element with article in html  
*title_selector* is css selector to find title in article  
*description_selector* is css selector to find description in article  
*link_selector* is css selector to find link in article or "self" if article is link