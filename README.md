# openbooks-cli
Openbooks is a command line utility to download ebooks from irchighway.net. 

## Features
- Search eBooks
  - Returns a file of search results to your `~/Downloads` folder
  - Example Results
  ```
  !LawdyServer Kurt Vonnegut - Slaughterhouse Five.pdf  ::INFO:: 220.0KB
  !Oatmeal Vonnegut, Kurt - Slaughterhouse-Five(v1.1)[htm].rar  ::INFO:: 257.0KB 
  !Oatmeal 023 - Vonnegut, Kurt Jr. - Slaughterhouse-five (v1.1) [html, jpg].rar  ::INFO:: 257.0KB 
  !Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v4.0) (html).rar  ::INFO:: 245.9KB 
  !Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v5.0) (epub).rar  ::INFO:: 1.9MB 
  !Pondering42 023 - Vonnegut, Kurt Jr. - Slaughterhouse-five (v1.1) [html, jpg].rar ::INFO:: 256.99KB
  ``` 
- Download eBooks
  - Enter the download string of the book you want to receive (don't include the ::INFO:: )
    - Ex) `!Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v5.0) (epub).rar`
  - The book file will be downloaded to your `~/Downloads` directory
  - If the book is contained within a compressed archive, it is extracted automatically

## How it works
- I wrote this program mainly as a way to learn the [IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) protocol. Many IRC clients also allow file downloads which are implemented using the [DCC](https://en.wikipedia.org/wiki/Direct_Client-to-Client) protocol. The essentials of both are implemented in order for the application to work.
