# Endpoints

/listings -> alle listings in der collection

/listings/{uuid} -> listing 



## JSON objects

### Listings

{
    uuid: uuid(string)
    Title: string,
    Description: string,
    Images: ? Binary Image (but has to support lists),
    Creator: ID,
    Tags: [strings],
    Paypal: Image (email),
    Text: string

}