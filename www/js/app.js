var id = "dashboard-"+(Math.random()+"").substr(2, 8)
var client = new Messaging.Client(location.hostname, location.port || 80, id);
client.onConnectionLost = onConnectionLost;
client.onMessageArrived = onMessageArrived;
client.connect({onSuccess:onConnect});

function onConnect() {
  console.log("onConnect");
  client.subscribe("devices/#");
  // message = new Messaging.Message("Hello");
  // message.destinationName = "/World";
  // client.send(message);
  // client.disconnect();
}

function onConnectionLost(responseObject) {
  if (responseObject.errorCode !== 0)
   console.log("onConnectionLost:", responseObject.errorMessage);
}
function onMessageArrived(message) {

  console.log(message);

  // console.log(message);
  if(window.view && window.view.onMessage)
    window.view.onMessage(message)

  /*
  var $target = $("[data-topic='"+message.destinationName+"']")
  if ($target[0]) {
    $target.append($("<li>").text(message.payloadString))
  }
  */
  // console.log("onMessageArrived:", message.payloadString);
}

$(window).on("hashchange", navigate)

navigate()

async function navigate() {
  if(window.view && window.view.onUnload)
    window.view.onUnload()
  window.view = null

  var container;
  var match = location.hash.match(/\#devices\/([\w\-._~:\[\]!$'\(\)*+,;=]+)$/);
  if (match) container = await inflateDevice(match[1]) //
  else {
    match = location.hash.match(/\#devices\/([\w\-._~:\[\]!$'\(\)*+,;=]+)\/sensors\/([\w\-._~:\[\]!$'\(\)*+,;=]+)$/);
    if (match) container = await inflateSensor(match[1], match[2])
    else {
      match = location.hash.match(/\#devices\/([\w\-._~:\[\]!$'\(\)*+,;=]+)\/actuators\/([\w:\-._~:\[\]!$'\(\)*+,;=]+)$/);
      if (match) container = await inflateActuator(match[1], match[2])
      else container = await inflateDevices()
    }
  }
  if (container) {
    $(document.body)
      .empty()
      .append(container)

    if(window.view && window.view.onLoad)
      window.view.onLoad()
  } else {
    alert("Navigation Error!")
  }
}

////////////////////////////////////////////////////////////////////////////////

async function inflateDevices() {

  var resp = await fetch("/devices")
  var devices = await resp.json();

  var $container = $("<div>")
    .addClass("devices container")
    .append(breadcrumb([
      {text: "devices", href: "devices"},
    ]));

  $ul = $("<ul>")
    .addClass("devices")
    .appendTo($container)

  for (var device of devices) {

    $("<li>")
      .addClass("device")
      .appendTo($ul)
      .append($("<a>")
        .addClass("id")
        .attr("href", "#devices/"+device.id)
        .append(device.name))
  }

  return $container[0]
}

////////////////////////////////////////////////////////////////////////////////

async function inflateDevice(deviceId) {

  var resp = await fetch("/devices/"+deviceId)
  var device = await resp.json();

  var $container = $("<div>")
    .addClass("device container")
    .append(breadcrumb([
      {text: "devices", href: "devices"},
      {text: deviceId, href: "devices/"+deviceId},
    ]));

  $("<h1>")
    .text(device.name)
    .appendTo($container)

  $("<h3>")
    .text("Sensors")
    .appendTo($container)

  $sensors = $("<ul>")
    .addClass("sensors")
    .appendTo($container);

  if (device.sensors) {
    for (var sensor of device.sensors) {
      $("<li>")
      .addClass("sensor")
      .appendTo($sensors)
      .append($("<a>")
      .attr("href", "#devices/"+device.id+"/sensors/"+sensor.id)
      .text(sensor.name))
    }
  }

  $("<hr>").appendTo($container)

  $("<h3>")
    .text("Actuators")
    .appendTo($container)

  $actuators = $("<ul>")
    .addClass("actuators")
    .appendTo($container);

  if (device.actuators) {
    for (var actuator of device.actuators) {
      $("<li>")
      .addClass("actuator")
      .appendTo($actuators)
      .append($("<a>")
      .attr("href", "#devices/"+device.id+"/actuators/"+actuator.id)
      .text(actuator.name))
    }
  }

  return $container[0]
}

////////////////////////////////////////////////////////////////////////////////

var colors = [
  'rgb(255, 99, 132)', 'rgb(255, 159, 64)', 'rgb(255, 205, 86)',
  'rgb(75, 192, 192)', 'rgb(54, 162, 235)', 'rgb(153, 102, 255)',
  'rgb(201, 203, 207)'];


async function inflateSensor(deviceId, sensorId) {

  var resp = await fetch("/devices/"+deviceId+"/sensors/"+sensorId)
  var sensor = await resp.json();

  var $container = $("<div>")
    .addClass("sensor container")
    .append(breadcrumb([
      {text: "devices", href: "devices"},
      {text: deviceId, href: "devices/"+deviceId},
      {text: "actuators", href: "devices/"+deviceId},
      {text: sensorId, href: "devices/"+deviceId+"/sensors/"+sensorId},
    ]));

  $("<h1>")
    .text(sensor.name)
    .appendTo($container)

  $("<ul>")
    .attr("data-topic", "devices/"+deviceId+"/sensors/"+sensorId+"/value")
    .appendTo($container)

  $canvas = $("<canvas>")
    .addClass("chart")
    .appendTo($container)

  var config = {
		type: 'line',
		data: {
			labels: [],
			datasets: []
		},
		options: {
			responsive: true,
			title: {
				display: true,
				text: sensor.name
			},
      tooltips: false,
			/*tooltips: {
				mode: 'index',
				intersect: false,
			},
			hover: {
				mode: 'nearest',
				intersect: true
			},*/
			scales: {
				xAxes: [{
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Time'
					}
				}],
				yAxes: [{
					display: true,
					scaleLabel: {
						display: false,
						labelString: ''
					}
				}]
			}
		}
	};

  var chart;
  var times = [];
  /*
  var now = new Date;
  for (var i=30; i>0; i--) {
    var time = new Date(now-i*1000)
    times.push(time*1)
    config.data.labels.push(time.getHours()+":"+time.getMinutes()+":"+time.getSeconds())
  }
  */

  window.view = {

    onMessage: (msg) => {

      console.log(msg);
      value = JSON.parse(msg.payloadString)
      if (msg.destinationName == "devices/"+deviceId+"/sensors/"+sensorId+"/values") {
        return
        values = value
        value = values[0]
        for(var i=1; i<values.length; i++) {
          for (var key in values[i]) {
            value[key] += values[i][key]
          }
        }
        for (var key in values) {
          value[key] /= values.length
        }
      } else if (msg.destinationName != "devices/"+deviceId+"/sensors/"+sensorId+"/value") {
        return
      }

      if (typeof value == "number")
        value = {[sensorId]: value}

      if (config.data.datasets.length == 0) {
        var i = 0;
        for (var key in value) {
          config.data.datasets.push({
    				label: key,
    				fill: false,
    				backgroundColor: colors[i],
    				borderColor: colors[i],
    				data: [value[key]],
    			});
          i++;
        }
      } else {
        var i = 0;
        for (var key in value) {
          config.data.datasets[i].data.push(value[key]);
          i++;
        }
      }
      /*
      var $target = $("[data-topic='"+message.destinationName+"']")
      if ($target[0]) {
        $target.append($("<li>").text(message.payloadString))
      }
      */

      var now = new Date()
      times.push(now*1)
      config.data.labels.push(now.getHours()+":"+now.getMinutes()+":"+now.getSeconds())

      var oldest = new Date(now*1-5*1000)*1
      while (oldest > times[0]) {
        times.shift()
        config.data.labels.shift()
        for (var dataset of config.data.datasets)
          dataset.data.shift()
      }

      chart.update();

    },
    onLoad: () => {

      var ctx = $canvas[0].getContext('2d');
  		chart = new Chart(ctx, config);
    }
  }

  return $container[0]
}

function randomScalingFactor() {
  return Math.round(Math.random()*200-100);
}

////////////////////////////////////////////////////////////////////////////////

async function inflateActuator(deviceId, actuatorId) {

  var resp = await fetch("/devices/"+deviceId+"/actuators/"+actuatorId)
  var actuator = await resp.json();

  var $container = $("<div>")
    .addClass("actuator container")
    .append(breadcrumb([
      {text: "devices", href: "devices"},
      {text: deviceId, href: "devices/"+deviceId},
      {text: "actuators", href: "devices/"+deviceId},
      {text: actuatorId, href: "devices/"+deviceId+"/actuators/"+actuatorId},
    ]));

  $("<h1>")
    .text(actuator.name)
    .appendTo($container)

  $("<p>")
    .appendTo($container)
    .text("New Value:")

  $textarea = $("<textarea>")
    .appendTo($container)
    .addClass("input")

  $submit = $("<button>")
    .appendTo($container)
    .text("Submit Value")
    .addClass("submit")
    .on("click", async () => {

      var resp = await fetch("/devices/"+deviceId+"/actuators/"+actuatorId+"/value", {
        method: "POST",
        body: $textarea.val(),
        headers: {
          "Content-Type": "application/json; charset=utf-8",
        }
      });
      var text = await resp.text();
      if (!resp.ok) alert("Error\n"+text)
      else alert("OK!\n"+text)
    });

  return $container[0]
}

////////////////////////////////////////////////////////////////////////////////

function breadcrumb(items) {

  var $breadcrumb = $("<ol>")
    .addClass("breadcrumb");

  for (var n=0; n<items.length; n++) {
    var item = items[n];

    $("<li>")
      .addClass("item")
      .appendTo($breadcrumb)
      .append($("<a>").attr("href", "#"+item.href).text(item.text));
  }

  return $breadcrumb[0]
}
