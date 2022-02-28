# OpenBooks API

> A self-hosted book search server.

I am getting frequent complaints about search spam from IRC admins. This API
will provide a service thats allows complaint free searching of irc.irchighway.net's books.

## How it works

1. Up-to-date search files are downloaded from IRC highway on a weekly basis.
   - These search files contain all the books each server has available
   - This involves making independent requests to each book server and downloading
     the master list of search result text files. For example typing "@Oatmeal"
     will send you all of Oatmeal's available books. Oatmeal's search file
     is 80MB with ~800,000 search strings.
2. Each line will be parsed into sections like author, title, server, etc.
3. Each line will be put into a full text search index.
4. The API will provide REST endpoints to enable book searches.
   - The time limits when using this service will be greatly reduced.
5. OpenBooks will by default send all search requests to this server.
   - Users can opt out via a command flag.
