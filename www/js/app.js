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
  var $target = $("[data-topic='"+message.destinationName+"']")
  if ($target[0]) {
    $target.append($("<li>").text(message.payloadString))
  }
  // console.log("onMessageArrived:", message.payloadString);
}

$(window).on("hashchange", navigate)

navigate()

async function navigate() {
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

  for (var sensor of device.sensors) {
    $("<li>")
      .addClass("sensor")
      .appendTo($sensors)
      .append($("<a>")
        .attr("href", "#devices/"+device.id+"/sensors/"+sensor.id)
        .text(sensor.name))
  }

  $("<hr>").appendTo($container)

  $("<h3>")
    .text("Actuators")
    .appendTo($container)

  $actuators = $("<ul>")
    .addClass("actuators")
    .appendTo($container);

  for (var actuator of device.actuators) {
    $("<li>")
      .addClass("actuator")
      .appendTo($actuators)
      .append($("<a>")
        .attr("href", "#devices/"+device.id+"/actuators/"+actuator.id)
        .text(actuator.name))
  }

  return $container[0]
}

////////////////////////////////////////////////////////////////////////////////

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

  return $container[0]
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
