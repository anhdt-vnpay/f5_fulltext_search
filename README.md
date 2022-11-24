# f5_fulltext_search

# start api (default port 8000)
make -f hanbq.yaml api 

# Demo model User
# API:

# GET: localhost:8000/es/get?table=...&condition=...

# BODY: json{table:"name", user: {"id", "name", "type"}}
# INSERT: localhost:8000/es/insert       => Post with BODY
# UPDATE: localhost:8000/es/update       => Post with BODY
# DELETE: localhost:8000/es/delete       => Post with BODY

# SEARCH LITE: localhost:8000/es/search?query="text_to_search"