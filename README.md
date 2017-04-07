Go Swagger "Lite"
===========================

This is a utility for generating API documentation from annotations in Go code. It generates the documentation as JSON, according to the [Swagger Spec](https://github.com/wordnik/swagger-spec) which can be used to view it in the [Swagger UI](https://github.com/wordnik/swagger-ui).

This library was originally forked from [yvasiyarov/swagger](https://github.com/yvasiyarov/swagger). That tool was inspired by [Beego](http://beego.me/docs/advantage/docs.md), and follows the same annotation standards set by Beego.

Like [yvasiyarov/swagger](https://github.com/yvasiyarov/swagger), this project is framework agnostic, but is simpler in an opinionated way.


Comment AnnotationsFormat
---------------------------

### 1. General API info

Use the following annotation comments to describe the API as a whole.
They should be placed in the "main" file of your application, above the "package" keyword.
The @-tags are not case sensitive (TODO Is this true?), but it is recommended to use the casing as shown, to be consistent.
Each of these annotations take a single argument that is an unquoted string to the end of the line.
The purpose of each annotation should be self-explanatory.
They are all optional, although using at least @Title and @Description is highly recommended.

    // @APIVersion 1.0.0
    // @Title My Cool API
    // @Description My API usually works as expected. But sometimes its not true
    // @Contact api@contact.me
    // @TermsOfServiceUrl http://google.com/
    // @License BSD
    // @LicenseUrl http://opensource.org/licenses/BSD-2-Clause



### 2. Sub API Definitions (One per Resource)

The Swagger specification is a bit confusing in how it refers to your API (singular) having multiple APIs (plural). It assumes that your API is "resource" centric. That is, it assumes that the first segment of every URL path refers to a "resource" and that it is therefore desirable to group the APIs specifications by these resources -- each with it's own API. (The [pet store example](http://petstore.swagger.wordnik.com/) is grouped according to three resources: pet, user, and store; with all of the URLs beginning with /pet/, /user/, or /store/, respectively).

NOTE: This is problematic for an application that is microservices-centric, rather than resource-centric. (See TODO, below.)

The @SubApi annotation is an opportunity to define each resource.

    // @SubApi Order management API [/order]
    // @SubApi Statistic gathering API [/cache-stats]

These @SubApi comments should also be placed above the "package" keyword of the "main" file of your application. You can declare several sub-API's, one after the other. The format of the @SubApi annotation is simple:

    // @SubApi DESCRIPTION [URI]

URI must have leading slash. The description is not mandatory, but if you forget it, then you will have an ugly looking document. :-)


### 3. API Operation

The most important annotation comments for swagger UI generation are these comments that describe an operation (GET, PUT, POST, etc.). These comments are placed within your controller source code. One full set of these comments is used for each operation. They are placed just above the code that handles that operation.

Please, refer to the following example when reviewing the notes, below:

    // @Title getOrdersByCustomer
    // @Description retrieves orders for given customer defined by customer ID
    // @Accept  json
    // @Param   customer_id     path    int     true        "Customer ID"
    // @Param   order_id        query   int     false        "Retrieve order with given ID only"
    // @Param   order_nr        query   string  false        "Retrieve order with given number only"
    // @Param   created_from    query   string  false        "Date-time string, MySQL format. If specified, API will retrieve orders that were created starting from created_from"
    // @Param   created_to      query   string  false        "Date-time string, MySQL format. If specified, API will retrieve orders that were created before created_to"
    // @Success 200 {array}  my_api.model.OrderRow
    // @Failure 400 {object} my_api.ErrorResponse    Customer ID must be specified
    // @Resource /order
    // @Router /orders/by-customer/{customer_id} [get]

Let's discuss every line in detail:
* The @Title provides a "nickname", in Swagger terms, to the operation. It is kind of an "alias" for this API operation. Only [A-Za-z0-9] characters are allowed. It's required, but only used internally. Swagger UI does not display it.
* @Description - A longer description for the operation. (An unquoted string to the end of line.)
* @Accept - One of: json, xml, plain, or html. (Can also be one of their fully qualified alternatives: application/json, text/xml, text/plain, or text/html.) Should be equal to the "Accept" header of your API.
* @Param - Defines a parameter that is accepted by this API operation. This comment has the following format:
 @Param  param_name  transport_type  data_type  required  "description"
 * param_name  - name of the parameter.
 * transport_type  - defines how this parameter is passed to the operation. Can be one of path/query/form/header/body
 * data_type  - type of parameter
 * required - Whether or not the parameter is mandatory (true or false).
 * description - parameter description. Must be quoted.
* @Success/@Failure - Use these annotations to define the possible responses by the API operation. The format is as follows:
 @Success http_response_code response_type response_data_type response_description
 * http_response_code 200 for success response, any other code for failure.
 * response_type - can be {object} or {array} -- depending on whether the operation returns a single JSON object, or an array of objects
 * response_data_type - data type of your response. Can be one of the Go built-in types, including error (which is actually a built-in interface, not a built-in type), or your custom type. All interface types except "error" will be displayed just as "interface". (It's not possible to find out which type it will has at parsing time.)
 * response_description - optional. It usually only makes sense for error responses.
* @Router - define route path, which should be used to call this API operation. It has the following format:
 @Router request_path [request_method]
 * reqiest_path which should be used by to make request to this API endpoint. It can include placeholders for parameters with transport type equal "path". Look at the example above.
 * request_method - just HTTP request method (get/post/etc..)
* @Resource - Forces the resource identifier to be something other than the first segment of the route URI. For example, if "@Resource /payment" is specified together with "@Router /invoice/{id}/payments [get]" then this operation will be part of the "payment" sub-api, rather than the "invoice" sub-api. This is also good for when the @Router specifies a path in which the first segment is almost correct, but not quite, e.g. /payments (plural) vs. /payment (singular). It has the following format:
@Resource resource_name
 * resource_name - A leading slash in the name in the @Resource annotation is optional. So, "@Resource /payment" and "@Resource payment" are the same.

### 4. Struct Tags

    type Actor struct {
        Id            *string   `required json:"id"`
        FirstName     *string   `json:"firstName,required" description:"The actor's first/middle name(s)"`
        LastName      *string   `json:"lastName,required" description:"The actor's last name (sorted by)"`
        HeadshotImage *string   `json:"-"`
        Filmography   []Film    `json:"filmography,omitempty"`
        Contact       *Contact  `json:"contact,omitempty"`
    }

* If a `required` struct tag is found, then the field is marked as required, e.g. `Id`, above.
* If `required` is found within a `json` struct tag, then the field is marked as required, e.g. `FirstName`, above.
* If a `description` struct tag is found, then it provides the field's description, e.g. `FirstName`, above.
* If `-` is found within a `json` struct tag, then the field is ignored (not documented), e.g. `HeadshotImage`, above.

Note: Use a space to separate multiple struct tags.


Quick Start Guide
-----------------

1. Add comments to your API source code.
2. Install `go get github.com/RobotsAndPencils/go-swaggerLite`
3. Run generator from your project:

        go-swaggerLite -apiPackage="github.com/myuser/myproject" \
            -mainApiFile="github.com/myuser/myproject/web/main.go" \
            -basePath="http://127.0.0.1:3000"

    Command line switches are:
    * -apiPackage  - package with API controllers implementation
    * -mainApiFile - main API file. We will look for "General API info" in this file. If the mainApiFile command-line switch is left blank, then main.go is assumed (in the location specified by apiPackage).
    * -basePath    - Your API URL. Test requests will be sent to this URL
    * -packageExclusionList - Comma delimited list of packages to report-continue vs. fail when not found

4. This will generate a `generatedSwaggerSpec.go` in `package main`. In this a `swaggerApiHandler` function is expossed that takes a URI path prefeix (to strip off) and returns a `http.HandlerFunc`. You might map this via a Gorilla Mux router like:

        r := mux.NewRouter()
        // ...
        r.HandleFunc("/spec", SwaggerApiHandler("/spec"))
        spec := r.PathPrefix("/spec").Subrouter()
        spec.HandleFunc("/{resource:.*}", SwaggerApiHandler("/spec"))
        // ...
        http.Handle("/", r)

5. Your Swagger API JSON description can be found out `<origin>/spec`.

Known Limitations
-----------------

* Interface types are not supported, because it's not possible to resolve them to actual implementations are parse-time. All interface values will be displayed just as "interface".
* Types that implement the Marshaler/Unmarshaler interface. Marshaling of this types will produce unpredictable JSON (at parse-time).
