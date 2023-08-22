load("render.star", "render")
load("time.star", "time")
load("cache.star", "cache")
load("http.star", "http")
load("encoding/json.star", "json")


DEFAULT_TEMP = 999

def main(config):
    
    print("STARTING.......")


    # setup fonts
    font = config.get("font", "tb-8")
    print("Using font: '{}'".format(font))

    #call api to get data before render
    baseURL = "http://localhost:1323/current"
    api_result = http.get(url = baseURL)
    api_response = api_result.body()
    cache.set("temps", api_result.body(), ttl_seconds = 7200)


    weather_list = json.decode(api_response)
    print(weather_list)

    # render screen
    return render.Root(
        render.Row(
            children=[
                render.Column(
                        expanded=True,
                        main_align="space_around",
                        cross_align="left",
                        children=[
                            render.Text("Inside:", color="#0096FF"),
                            render.Text("Outside:", color="#0047AB"),
                            render.Text("Rain:", color="#3F00FF"),
                        ]
                ),
                render.Column(
                        expanded=True,
                        main_align="space_around",
                        cross_align="left",
                        children=[
                            render.Text(str(weather_list["insideTemp"]), color="#7393B3"),
                            render.Text(str(weather_list["outsideTemp"]), color="#088F8F"),
                            render.Text(str(weather_list["rainfall"]), color="#6495ED"),
                        ]
                )
            ]
        )
        
)
    

    
    