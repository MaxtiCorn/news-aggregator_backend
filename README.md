# NewsAggregator

## Build

Run `go get` to get dependencies (may need [gcc](http://tdm-gcc.tdragon.net/download) for this)
Run `go build` to build the project, it will generate an exe file

## Run

*exe_name* -config="*path_to_config_file*" -db="path_to_database_file"
where path_to_config_file is path to json file with parsing rules (default: config.json located in the same directory as the exe file),
      path_to_database_file is path to database (default: database.db located in the same directory as the exe file)

Config schema:
{
    "rss": [                                    *list of parsing rules for rss*
        {
            "url": string,                      *url of rss*
            "interval": number                  *polling interval in seconds*
        }
    ],
    "html": [                                   *list of parsing rules for sites*
        {
            "url": string,                      *url of site*
            "interval": number,                 *polling interval in seconds*
            "article_selector": string,         *selector for find article* 
            "title_selector": string,           *selector for find title in article*  
            "description_selector": string      *selector for find description in article*
        }
    ]
}