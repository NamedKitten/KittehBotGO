$(function () {
  window.dps = [];
  window.xVal = 1;
  window.yVal = 0;
  window.updateInterval = 0;

  $.ajax({
    type: 'GET',
    url: '/interval',
    complete: function (responseText) {
      console.log(responseText)
      window.updateInterval = parseFloat(responseText.responseText)
    }
  });

  window.memdps = [{ x: 0, y: 0 }];
  var socket = io();

  CanvasJS.addColorSet("customColorSet", [
    "#393f63",
    "#e5d8B0",
    "#ffb367",
    "#f98461",
    "#d9695f",
    "#e05850"
  ]);

  window.memUsage = new CanvasJS.Chart("memUsage", {
    animationEnabled: true,
    animation: true,
    responsive: true,
    backgroundColor: "transparent",
    axisX: {
      labelFontColor: "#717171",
      lineColor: "#a2a2a2",
      tickColor: "transparent",
      title: "Time (seconds?)"
      
    },
    axisY: {
      gridThickness: 0,
      labelFontColor: "#717171",
      lineColor: "#a2a2a2",
      tickColor: "#a2a2a2",
      minimum: 2,
      title: "Memory Used (MB)"
    },
    toolTip: { cornerRadius: 0, fontStyle: "normal" },
    data: [{ color: "#CD5740", type: "line", dataPoints: dps }]
  });
  window.memUsage.render();
  window.memStats = new CanvasJS.Chart("memStats", {
    animationDuration: window.updateInterval,
    animationEnabled: true,
    backgroundColor: "transparent",
    colorSet: "customColorSet",
    axisX: {
      labelFontColor: "#717171",
      labelFontSize: 18,
      lineThickness: 0,
      tickThickness: 0,
      lineColor: "#a2a2a2",
      tickColor: "transparent"
    },
    axisY: {
      gridThickness: 0,
      labelFontColor: "#717171",
      lineColor: "#a2a2a2",
      tickColor: "#a2a2a2",
      minimum: 0,
      title: "Memory (MB)"
    },
    data: [
      {
        indexLabelFontColor: "#717171",
        indexLabelFontFamily: "calibri",
        indexLabelFontSize: 18,
        indexLabelPlacement: "outside",
        indexLabelFormatter: function (e) {
          return e.dataPoint.y.toFixed(5) + " MB";
        },
        type: "bar",
        dataPoints: window.memdps
      }
    ]
  });
  window.memStats.render();
  /*
  var updateChart = function() {
    $.ajax({
      url: "getdata",
      type: "GET",
      success: function(data) {()*/
  window.setTimeout(function () {

    socket.on('mem stats', function (msg) {
      console.log("stuff!")
      var stuff = msg.split("\n");
      window.yVal = parseFloat(stuff[0]);
      dps.push({ x: xVal, y: yVal });
      if (window.xVal > 20) {
        window.dps.shift();
      }
      window.memdps.splice(0, memdps.length);
      window.memdps.push(
        { y: parseFloat(stuff[0]), label: "Using" },
        { y: parseFloat(stuff[1]), label: "Allocated" },
        { y: parseFloat(stuff[2]), label: "Cleaned" }
      );
      window.memUsage.render();
      window.memStats.render();

      window.xVal += window.updateInterval / 1000;
    });
  }, 1000);
  //}
  /*});
};*/
  /*
  setInterval(function() {
    socket.emit('mem get', "get");
  }, window.updateInterval);
  */
});
