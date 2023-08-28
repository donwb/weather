load("render.star", "render")
load("time.star", "time")
load("cache.star", "cache")
load("http.star", "http")
load("encoding/json.star", "json")


DEFAULT_TEMP = 999

def main(config):
    
    print("STARTING.......")


    # setup fonts
    font = config.get("font", "5x8")
    print("Using font: '{}'".format(font))

    #call api to get data before render
    # baseURL = "http://localhost:1323/current"
    baseURL = "https://weathersrv-hf9df.ondigitalocean.app/current"
    api_result = http.get(url = baseURL)
    api_response = api_result.body()
    cache.set("temps", api_result.body(), ttl_seconds = 7200)


    weather_list = json.decode(api_response)
    print(weather_list)

    insideTempColor, outsideColor, humidityColor = computeColors(weather_list["insideTemp"], weather_list["outsideTemp"], weather_list["humidity"])

    # render screen
    return render.Root(
        render.Row(
            children=[
                render.Column(
                        expanded=True,
                        main_align="space_around",
                        cross_align="left",
                        children=[
                            render.Text("Inside:"),
                            render.Text("Outside:"),
                            render.Text("Humidity:"),
                            render.Text("Rain:"),
                        ]
                ),
                render.Column(
                        expanded=True,
                        children=[
                            render.Box(
                                width=4,
                                height=10,)
                        ],
                ),
                render.Column(
                        expanded=True,
                        main_align="space_around",
                        cross_align="left",
                        children=[
                            render.Text(str(weather_list["insideTemp"]), color=insideTempColor),
                            render.Text(str(weather_list["outsideTemp"]), color=outsideColor),
                            render.Text(str(weather_list["humidity"]), color=humidityColor),
                            render.Text(str(weather_list["rainfall"]), color="#6495ED"),
                        ]
                )
            ]
        )
        
)
    

def computeColors(inside, outside, humidity):
    red = "#B81D13"
    yellow = "EFB700"
    green = "008450"

    if inside < 70:
        insideColor = green
    elif inside > 70 and inside < 72:
        insideColor = yellow
    else:
        insideColor = red

    if outside < 80:
        outsideColor = green
    elif outside > 80 and outside < 90:
        print(outside)
        outsideColor = yellow
    else:
        outsideColor = red
    
    if humidity < 60:
        humidityColor = green
    elif humidity > 60 and humidity < 76:
        humidityColor = yellow
    else:
        humidityColor = red

    print(insideColor, outsideColor, humidityColor)
    return insideColor, outsideColor, humidityColor
