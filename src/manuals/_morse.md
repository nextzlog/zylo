### zylo/morse

{% capture body %}
{{.EmitUsage}}
{% endcapture %}
{{`{{body | replace: "## Usage", "" | markdownify}}`}}
