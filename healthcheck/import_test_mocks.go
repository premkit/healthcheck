package healthcheck

const (
	// This is formatted carefully and correctly for YAML.
	simpleServiceYAML = `- name: service1
  checks:
    - type: http
      spec:
        method: "GET"
        options: ['no_follow_redirects', 'ignore_insecure']
        status_codes_available: [200]
        uri: "http://localhost"
`

	// The above service, translated to JSON.
	simpleServiceJSON = `[
		{
			"name": "service1",
			"checks": [
				{
					"type": "http",
					"spec": {
						"method": "POST",
						"options": [
							"no_follow_redirects",
							"ignore_insecure"
						],
						"status_codes_available": [
							200
						],
						"uri": "http://localhost"
					}
				}
			]
		}
	]`
)
