function addServersHTML(guilds) {
  var htmlString = `

    ${guilds.map(guild => `	<li><div class="media">
    <img src="${guild.icon}" height="64px" width="64px" class=mr-3>
    <div class="media-body">
      <h5 class="mt-0">${guild.name}</h5>
      ${guild.members} members
    </div></li>`)}
`;
  document.getElementById("guilds-list").innerHTML = htmlString;
}

function addPlayersHTML(players) {
  var htmlString = `

    ${players.map(player => `	<li><div class="media">
    <img src="${player.thumbnail}" height="64px" width="64px" class=mr-3>
    <div class="media-body">
      <h5 class="mt-0">${player.title}</h5>
      Playing in ${player.guildName}
    </div></li>`)}
`;
  document.getElementById("music-players-list").innerHTML = htmlString;
}

function compareGuilds(a, b) {
  if (a.members < b.members) return 1;
  if (b.members < a.members) return -1;

  return 0;
}

document.addEventListener('DOMContentLoaded', function() {
  window.dps = [];

  window.memdps = [{
      x: 0,
      y: 0
  }];
  var socket = io();
  window.connectionTime = 0;

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
      toolTip: {
          cornerRadius: 0,
          fontStyle: "normal"
      },
      data: [{
          color: "#CD5740",
          type: "line",
          dataPoints: dps
      }]
  });
  window.memUsage.render();
  window.memStats = new CanvasJS.Chart("memStats", {
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
      data: [{
          indexLabelFontColor: "#717171",
          indexLabelFontFamily: "Lato",
          indexLabelFontSize: 18,
          indexLabelPlacement: "outside",
          indexLabelFormatter: function(e) {
              return e.dataPoint.y.toFixed(5) + " MB";
          },
          type: "bar",
          dataPoints: window.memdps
      }]
  });
  window.memStats.render();
  socket.on('connect', function(msg) {
      console.log("connected")
      window.connectionTime = Date.now()
  });

  socket.on('mem stats', function(msg) {
      var jsonMsg = JSON.parse(msg)
      timeElapsed = (Date.now() - window.connectionTime) / 1000;

      dps.push({
          x: timeElapsed,
          y: jsonMsg["using"]
      });
      if (timeElapsed > 20) {
          window.dps.shift();
      }
      window.memdps.splice(0, memdps.length);
      window.memdps.push({
          y: jsonMsg["using"],
          label: "Using"
      }, {
          y: jsonMsg["allocated"],
          label: "Allocated"
      }, {
          y: jsonMsg["cleaned"],
          label: "Cleaned"
      });
      window.memUsage.render();
      window.memStats.render();

  });

  socket.on('guilds stats', function(msg) {
      var jsonMsg = JSON.parse(msg)
      if (typeof jsonMsg === 'undefined' || jsonMsg === null) {
        return
      }
      jsonMsg.sort(compareGuilds);
      addServersHTML(jsonMsg)
  });

  socket.on('music stats', function(msg) {
      var jsonMsg = JSON.parse(msg)
      if (typeof jsonMsg === 'undefined' || jsonMsg === null) {
          return
      }
      addPlayersHTML(jsonMsg)
  });

  socket.on('disconnect', function(msg) {
      window.dps = [];
  });

});