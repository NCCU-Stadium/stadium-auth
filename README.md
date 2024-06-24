# auth-service

- RESTful APIs running on `localhost:8085`
- gRPC API running on `localhost:50050`

---

## RESTful APIs

### Login

<details>
<summary><code>POST</code> <code><b>/login/google</b></code> <code>(Login with Google)</code></summary>

##### Body (application/json or application/x-www-form-urlencoded)

> | key     | required | data type | description          |
> | ------- | -------- | --------- | -------------------- |
> | idToken | true     | string    | Id token from Google |

##### Responses

> | http code    | content-type       | response                                                                                                      |
> | ------------ | ------------------ | ------------------------------------------------------------------------------------------------------------- |
> | `200`        | `application/json` | `{"name": "username", "email": "user@ex.com", "avatar": "http://avatar.com", "refreshToken": "refreshToken"}` |
> | `401`, `500` | `text/plain`       | N/A                                                                                                           |

##### Cookie Setup

> Upon successful login, a cookie named `accessToken` will be set with HttpOnly and Secure flag.

</details>

<details>
<summary><code>POST</code> <code><b>/login/credentials</b></code> <code>(Login with credentials)</code></summary>

##### Body (application/json or application/x-www-form-urlencoded)

> | key      | required | data type | description |
> | -------- | -------- | --------- | ----------- |
> | email    | true     | string    | N/A         |
> | password | true     | string    | N/A         |

##### Responses

> | http code    | content-type       | response                                                                                                      |
> | ------------ | ------------------ | ------------------------------------------------------------------------------------------------------------- |
> | `200`        | `application/json` | `{"name": "username", "email": "user@ex.com", "avatar": "http://avatar.com", "refreshToken": "refreshToken"}` |
> | `401`, `500` | `text/plain`       | N/A                                                                                                           |

##### Cookie Setup

> Upon successful login, a cookie named `accessToken` will be set with HttpOnly and Secure flag.

</details>

### Register

<details>
<summary><code>POST</code> <code><b>/register</b></code> <code>(Register with mail and password)</code></summary>

##### Body (application/json or application/x-www-form-urlencoded)

> | key      | required | data type | description        |
> | -------- | -------- | --------- | ------------------ |
> | name     | true     | string    | N/A                |
> | email    | true     | string    | N/A                |
> | password | true     | string    | N/A                |
> | avatar   | false    | string    | Avatar's image URL |

##### Responses

> | http code    | content-type | response |
> | ------------ | ------------ | -------- |
> | `200`        | `text/plain` | N/A      |
> | `401`, `500` | `text/plain` | N/A      |

</details>

### Refresh Token

<details>
<summary><code>POST</code> <code><b>/refresh</b></code> <code>(Generate a new refresh token)</code></summary>

##### Body (application/json or application/x-www-form-urlencoded)

> | key          | required | data type | description |
> | ------------ | -------- | --------- | ----------- |
> | refreshToken | true     | string    | N/A         |

##### Responses

> | http code    | content-type       | response                                                                                                      |
> | ------------ | ------------------ | ------------------------------------------------------------------------------------------------------------- |
> | `200`        | `application/json` | `{"name": "username", "email": "user@ex.com", "avatar": "http://avatar.com", "refreshToken": "refreshToken"}` |
> | `401`, `500` | `text/plain`       | N/A                                                                                                           |

##### Cookie Setup

> Upon successful refresh, a cookie named `accessToken` will be set with HttpOnly and Secure flag.

</details>

## gRPC APIs

The auth-service also expose the functionality the generate and verify the JWT tokens.

### Verify

<details>
<summary><code>GET</code> <code><b>/verify</b></code> <code>(Verify a existing JWT token)</code></summary>

##### Request Type

> | key   | required | data type | description                 |
> | ----- | -------- | --------- | --------------------------- |
> | token | true     | string    | Starts with `Bearer<space>` |

##### Response Type

> | key     | data type | description                                          |
> | ------- | --------- | ---------------------------------------------------- |
> | claims  | string    | Stringified JSON data (Empty string if it's expired) |
> | expired | bool      | N/A                                                  |

##### Error Code

> | status code                | response       |
> | -------------------------- | -------------- |
> | `UNKNOWN 2`, `INTERNAL 13` | Invalid token  |
> | `INTERNAL 13`              | Marchal failed |

</details>
