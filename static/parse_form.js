function sendJSON(json) {
    var request = new XMLHttpRequest();
    var url = "/api/email"
    request.open("POST", url, true);
    request.setRequestHeader("Content-type",
			     "application/json");
    request.send(json);
    var log = "json sent:\n".concat(json);
    console.log(log);
}

(function () {
    function toJSONString(form) {
	var map = {};
	var elements = form.querySelectorAll(
	    "input, select, textarea");

	for(var i = 0; i < elements.length; i++) {
	    var e = elements[i];
	    var name = e.name;
	    var value = e.value;

	    if(name) {
		map[name] = value;
	    }
	}

	return JSON.stringify(map);
    }

    document.addEventListener(
	"DOMContentLoaded",
	function() {
	    var form = document.getElementById("email-form");
	    form.addEventListener("submit", function(e) {
		e.preventDefault();
		var json = toJSONString(this);
		console.log(json);
		sendJSON(json)
	    }, false);
	})
})();
