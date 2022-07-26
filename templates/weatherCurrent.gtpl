🌡️ Текущая погода в {{ .Name }}

{{ .GetIcon }} Сейчас: {{ .GetTemp }}, ощущается как {{ .GetFeelsLike }},
{{ (index .Weather 0).Description }}.

{{ if .Rain.OneHour }}
За последний час выпало {{ .Rain.OneHour }} мм дождя
{{ if .Rain.ThreeHours }}{{ .Rain.ThreeHours }} мм за последние 3 часа.{{ end }}
{{ else if .Snow.OneHour }}
За последний час выпало {{ .Rain.OneHour }} мм снега
{{ if .Snow.ThreeHours }}{{ .Snow.ThreeHours }} мм за последние 3 часа.{{ end }}
{{ else }}
За последние 3 часа осадков не выпало.
{{ end }}

Ветер: {{ .GetWindDirection }} {{ .Wind.Speed }} м/с, влажность: {{ .Main.Humidity }}%
Атмосферное давление: {{ .GetPressure | printf "%.1f"}} мм.рт.ст.
Видимость: {{ .Visibility }} м
Облачность: {{ .Clouds.All }}%