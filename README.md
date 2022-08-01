# go2ts

Easy create Typescript api from go struct or func

## go struct to typescript interface

```go
type Hello struct {
	Name string `json:"name" validate:"required"`
	Age  string `json:"age"`
	Vip  bool   `json:"vip,omitempty"`
}

type World struct {
	Dog  string `json:"dog" ts_type:"any"`
	Fish string `json:"fish,omitempty" validate:"required"`
}

func main() {
  go2ts.New().Add(Hello{}).Add(World{}).Write("hello.ts")
}
```

out put:

```ts
/* eslint-disable */

export interface Hello {
  name: string;
  age?: string;
  vip?: boolean;
}
export interface World {
  dog?: any;
  fish?: string;
}
```

## Auto create fetch API to Typescript

```go
type Hello struct {
	Name string `json:"name" validate:"required"`
	Age  string `json:"age"`
	Vip  bool   `json:"vip,omitempty"`
}

type World struct {
	Dog  string `json:"dog" ts_type:"any"`
	Fish string `json:"fish,omitempty" validate:"required"`
}

func main(){
  	go2ts.New().AddApi("POST", "/v1/world", GetWord).Write("apis.ts")
}
```

output:

```ts
// Auto create with go2ts
/* eslint-disable */

export const apiGetWord = (hello: Hello): Promise<World> => {
  return (window as any).customFetch("POST", "/v1/world", hello);
};
export interface Hello {
  name: string;
  age?: string;
  vip?: boolean;
}
export interface World {
  dog?: any;
  fish?: string;
}
```

You can create customFetch in window, like:

```js
window.customFetch = (method, url, body) => {
  return fetch(url, { method, body: JSON.stringify(body) }).then((v) =>
    v.json()
  );
};
```

And in your submit use the api, like:

```jsx
export default function () {
  const [res, setRes] = useState({ dog: "", fish: "" });
  useEffetc(() => {
    apiGetWord({
      name: "dog",
    }).then(setRes);
  }, []);

  return (
    <div>
      {res.dog} {res.fish}
    </div>
  );
}
```
