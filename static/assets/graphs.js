$(function () {

	CanvasJS.addColorSet("customColorSet", [
		"#393f63",
		"#e5d8B0",
		"#ffb367",
		"#f98461",
		"#d9695f",
		"#e05850",
	]);


  var dps = [];   //dataPoints.

        var xVal = 0;
        var yVal = 0;
        var updateInterval = 1000;
        var memdps = [];

	var memUsage = new CanvasJS.Chart("memUsage", {
		animationDuration: 1000,
		animationEnabled: true,
		backgroundColor: "transparent",
		axisX: {
			labelFontColor: "#717171",
			lineColor: "#a2a2a2",
			tickColor: "transparent",
		},
		axisY: {
			gridThickness: 0,
			labelFontColor: "#717171",
			lineColor: "#a2a2a2",
			tickColor: "#a2a2a2",
      minimum: 1.5

		},
		toolTip: {
			cornerRadius: 0,
			fontStyle: "normal",
		},
		data: [
			{
				color: "#CD5740",
				type: "line",
				dataPoints : dps
			}
		]
	});

  memUsage.render();
  var memStats = new CanvasJS.Chart("memStats", {
    animationDuration: 800,
    animationEnabled: true,
    backgroundColor: "transparent",
    colorSet: "customColorSet",
    axisX: {
      labelFontColor: "#717171",
      labelFontSize: 18,
      lineThickness: 0,
      tickThickness: 0,
      lineColor: "#a2a2a2",
      tickColor: "transparent",
    },
    axisY: {
			gridThickness: 0,
			labelFontColor: "#717171",
			lineColor: "#a2a2a2",
			tickColor: "#a2a2a2",
      minimum: 1.5

		},
    data: [
      {
        indexLabelFontColor: "#717171",
        indexLabelFontFamily: "calibri",
        indexLabelFontSize: 18,
        indexLabelPlacement: "outside",
        indexLabelFormatter: function (e) {
          return  e.dataPoint.y.toFixed(5);
        },
        type: "bar",
        dataPoints: memdps
      }
    ]
  });
  memStats.render();



  var updateChart = function () {

    $.ajax({
//The URL to process the request
'url' : 'getdata',
//The type of request, also known as the "method" in HTML forms
//Can be 'GET' or 'POST'
'type' : 'GET',
'success' : function(data) {
//You can use any jQuery/JavaScript here!!!
var stuff = data.split("\n")

yVal =  parseFloat(stuff[0])
dps.push({x: xVal,y: yVal});
//memdps = []
memdps.splice(0,memdps.length)

memdps.push({y: parseFloat(stuff[0]), label: "Using"}, {y: parseFloat(stuff[1]), label: "Allocated"}, {y: parseFloat(stuff[2]), label: "Cleaned"})
//console.log(memdps)
memUsage.render();
memStats.render();
xVal++;
}
});
}



  setInterval(function(){updateChart()}, 1000);


});
