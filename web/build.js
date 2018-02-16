(function(){function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s}return e})()({1:[function(require,module,exports){
"use strict";

var _other = require("./other");

console.log((0, _other.Foo)());

fetch("/data.json").then(function (resp) {
  resp.json().then(function (data) {
    console.log(data);
    var dat = JSON.stringify(data);

    var rows = [];
    var header = [React.createElement("th", { key: "empty" })];
    for (var j = 0; j < data["RunIDs"].length; j++) {
      var rid = data["RunIDs"][j];
      var run = data["Runs"][rid];
      header.push(React.createElement(
        "th",
        { key: "run-th-" + rid },
        run.Name
      ));
    }

    for (var i = 0; i < data["WorkflowIDs"].length; i++) {
      var wfid = data["WorkflowIDs"][i];
      var wf = data["Workflows"][wfid];
      var cells = [];
      cells.push(React.createElement(
        "td",
        { key: "wf-name-" + wfid },
        wf.Name
      ));
      for (var j = 0; j < data["RunIDs"].length; j++) {
        var rid = data["RunIDs"][j];
        var run = data["Runs"][rid];
        var classname = "";
        if (run.Total == run.Complete) {
          classname = "complete";
        }
        cells.push(React.createElement(
          "td",
          { key: "run-" + rid, className: classname },
          run.Complete,
          " / ",
          run.Total
        ));
      }
      rows.push(React.createElement(
        "tr",
        { key: "wf-" + wfid },
        cells
      ));
    }

    ReactDOM.render(React.createElement(
      "div",
      null,
      React.createElement(
        "h3",
        null,
        "Mortar"
      ),
      React.createElement(
        "table",
        null,
        React.createElement(
          "thead",
          null,
          React.createElement(
            "tr",
            null,
            header
          )
        ),
        React.createElement(
          "tbody",
          null,
          rows
        )
      )
    ), document.getElementById('root'));
  });
});

ReactDOM.render(React.createElement(
  "div",
  null,
  React.createElement(
    "h3",
    null,
    "Mortar"
  ),
  React.createElement(
    "div",
    null,
    "Loading"
  )
), document.getElementById('root'));

},{"./other":2}],2:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
function Foo() {
  return "Foo";
}

exports.Foo = Foo;

},{}]},{},[1])
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIm5vZGVfbW9kdWxlcy9icm93c2VyLXBhY2svX3ByZWx1ZGUuanMiLCJpbmRleC5qcyIsIm90aGVyLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBOzs7QUNBQTs7QUFFQSxRQUFRLEdBQVIsQ0FBWSxpQkFBWjs7QUFFQSxNQUFNLFlBQU4sRUFBb0IsSUFBcEIsQ0FBeUIsVUFBUyxJQUFULEVBQWU7QUFDdEMsT0FBSyxJQUFMLEdBQVksSUFBWixDQUFpQixVQUFTLElBQVQsRUFBZTtBQUM5QixZQUFRLEdBQVIsQ0FBWSxJQUFaO0FBQ0EsUUFBSSxNQUFNLEtBQUssU0FBTCxDQUFlLElBQWYsQ0FBVjs7QUFFQSxRQUFJLE9BQU8sRUFBWDtBQUNBLFFBQUksU0FBUyxDQUFDLDRCQUFJLEtBQUksT0FBUixHQUFELENBQWI7QUFDQSxTQUFLLElBQUksSUFBSSxDQUFiLEVBQWdCLElBQUksS0FBSyxRQUFMLEVBQWUsTUFBbkMsRUFBMkMsR0FBM0MsRUFBZ0Q7QUFDOUMsVUFBSSxNQUFNLEtBQUssUUFBTCxFQUFlLENBQWYsQ0FBVjtBQUNBLFVBQUksTUFBTSxLQUFLLE1BQUwsRUFBYSxHQUFiLENBQVY7QUFDQSxhQUFPLElBQVAsQ0FBWTtBQUFBO0FBQUEsVUFBSSxLQUFLLFlBQVksR0FBckI7QUFBMkIsWUFBSTtBQUEvQixPQUFaO0FBQ0Q7O0FBRUQsU0FBSyxJQUFJLElBQUksQ0FBYixFQUFnQixJQUFJLEtBQUssYUFBTCxFQUFvQixNQUF4QyxFQUFnRCxHQUFoRCxFQUFxRDtBQUNuRCxVQUFJLE9BQU8sS0FBSyxhQUFMLEVBQW9CLENBQXBCLENBQVg7QUFDQSxVQUFJLEtBQUssS0FBSyxXQUFMLEVBQWtCLElBQWxCLENBQVQ7QUFDQSxVQUFJLFFBQVEsRUFBWjtBQUNBLFlBQU0sSUFBTixDQUFXO0FBQUE7QUFBQSxVQUFJLEtBQUssYUFBYSxJQUF0QjtBQUE2QixXQUFHO0FBQWhDLE9BQVg7QUFDQSxXQUFLLElBQUksSUFBSSxDQUFiLEVBQWdCLElBQUksS0FBSyxRQUFMLEVBQWUsTUFBbkMsRUFBMkMsR0FBM0MsRUFBZ0Q7QUFDOUMsWUFBSSxNQUFNLEtBQUssUUFBTCxFQUFlLENBQWYsQ0FBVjtBQUNBLFlBQUksTUFBTSxLQUFLLE1BQUwsRUFBYSxHQUFiLENBQVY7QUFDQSxZQUFJLFlBQVksRUFBaEI7QUFDQSxZQUFJLElBQUksS0FBSixJQUFhLElBQUksUUFBckIsRUFBK0I7QUFDN0Isc0JBQVksVUFBWjtBQUNEO0FBQ0QsY0FBTSxJQUFOLENBQVc7QUFBQTtBQUFBLFlBQUksS0FBSyxTQUFTLEdBQWxCLEVBQXVCLFdBQVcsU0FBbEM7QUFBOEMsY0FBSSxRQUFsRDtBQUFBO0FBQStELGNBQUk7QUFBbkUsU0FBWDtBQUNEO0FBQ0QsV0FBSyxJQUFMLENBQVU7QUFBQTtBQUFBLFVBQUksS0FBSyxRQUFRLElBQWpCO0FBQXdCO0FBQXhCLE9BQVY7QUFDRDs7QUFFRCxhQUFTLE1BQVQsQ0FDRztBQUFBO0FBQUE7QUFDQztBQUFBO0FBQUE7QUFBQTtBQUFBLE9BREQ7QUFFQztBQUFBO0FBQUE7QUFDRTtBQUFBO0FBQUE7QUFBTztBQUFBO0FBQUE7QUFBSztBQUFMO0FBQVAsU0FERjtBQUVFO0FBQUE7QUFBQTtBQUFRO0FBQVI7QUFGRjtBQUZELEtBREgsRUFRRSxTQUFTLGNBQVQsQ0FBd0IsTUFBeEIsQ0FSRjtBQVVILEdBdkNDO0FBdUNDLENBeENIOztBQTBDQSxTQUFTLE1BQVQsQ0FDRztBQUFBO0FBQUE7QUFDQztBQUFBO0FBQUE7QUFBQTtBQUFBLEdBREQ7QUFDZ0I7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQURoQixDQURILEVBSUUsU0FBUyxjQUFULENBQXdCLE1BQXhCLENBSkY7Ozs7Ozs7O0FDOUNBLFNBQVMsR0FBVCxHQUFlO0FBQ2IsU0FBTyxLQUFQO0FBQ0Q7O1FBRVEsRyxHQUFBLEciLCJmaWxlIjoiZ2VuZXJhdGVkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXNDb250ZW50IjpbIihmdW5jdGlvbigpe2Z1bmN0aW9uIGUodCxuLHIpe2Z1bmN0aW9uIHMobyx1KXtpZighbltvXSl7aWYoIXRbb10pe3ZhciBhPXR5cGVvZiByZXF1aXJlPT1cImZ1bmN0aW9uXCImJnJlcXVpcmU7aWYoIXUmJmEpcmV0dXJuIGEobywhMCk7aWYoaSlyZXR1cm4gaShvLCEwKTt2YXIgZj1uZXcgRXJyb3IoXCJDYW5ub3QgZmluZCBtb2R1bGUgJ1wiK28rXCInXCIpO3Rocm93IGYuY29kZT1cIk1PRFVMRV9OT1RfRk9VTkRcIixmfXZhciBsPW5bb109e2V4cG9ydHM6e319O3Rbb11bMF0uY2FsbChsLmV4cG9ydHMsZnVuY3Rpb24oZSl7dmFyIG49dFtvXVsxXVtlXTtyZXR1cm4gcyhuP246ZSl9LGwsbC5leHBvcnRzLGUsdCxuLHIpfXJldHVybiBuW29dLmV4cG9ydHN9dmFyIGk9dHlwZW9mIHJlcXVpcmU9PVwiZnVuY3Rpb25cIiYmcmVxdWlyZTtmb3IodmFyIG89MDtvPHIubGVuZ3RoO28rKylzKHJbb10pO3JldHVybiBzfXJldHVybiBlfSkoKSIsImltcG9ydCB7IEZvbyB9IGZyb20gJy4vb3RoZXInXG5cbmNvbnNvbGUubG9nKEZvbygpKVxuXG5mZXRjaChcIi9kYXRhLmpzb25cIikudGhlbihmdW5jdGlvbihyZXNwKSB7XG4gIHJlc3AuanNvbigpLnRoZW4oZnVuY3Rpb24oZGF0YSkge1xuICAgIGNvbnNvbGUubG9nKGRhdGEpXG4gICAgdmFyIGRhdCA9IEpTT04uc3RyaW5naWZ5KGRhdGEpO1xuXG4gICAgdmFyIHJvd3MgPSBbXTtcbiAgICB2YXIgaGVhZGVyID0gWzx0aCBrZXk9XCJlbXB0eVwiPjwvdGg+XVxuICAgIGZvciAodmFyIGogPSAwOyBqIDwgZGF0YVtcIlJ1bklEc1wiXS5sZW5ndGg7IGorKykge1xuICAgICAgdmFyIHJpZCA9IGRhdGFbXCJSdW5JRHNcIl1bal07XG4gICAgICB2YXIgcnVuID0gZGF0YVtcIlJ1bnNcIl1bcmlkXTtcbiAgICAgIGhlYWRlci5wdXNoKDx0aCBrZXk9e1wicnVuLXRoLVwiICsgcmlkfT57cnVuLk5hbWV9PC90aD4pO1xuICAgIH1cblxuICAgIGZvciAodmFyIGkgPSAwOyBpIDwgZGF0YVtcIldvcmtmbG93SURzXCJdLmxlbmd0aDsgaSsrKSB7XG4gICAgICB2YXIgd2ZpZCA9IGRhdGFbXCJXb3JrZmxvd0lEc1wiXVtpXTtcbiAgICAgIHZhciB3ZiA9IGRhdGFbXCJXb3JrZmxvd3NcIl1bd2ZpZF07XG4gICAgICB2YXIgY2VsbHMgPSBbXTtcbiAgICAgIGNlbGxzLnB1c2goPHRkIGtleT17XCJ3Zi1uYW1lLVwiICsgd2ZpZH0+e3dmLk5hbWV9PC90ZD4pO1xuICAgICAgZm9yICh2YXIgaiA9IDA7IGogPCBkYXRhW1wiUnVuSURzXCJdLmxlbmd0aDsgaisrKSB7XG4gICAgICAgIHZhciByaWQgPSBkYXRhW1wiUnVuSURzXCJdW2pdO1xuICAgICAgICB2YXIgcnVuID0gZGF0YVtcIlJ1bnNcIl1bcmlkXTtcbiAgICAgICAgdmFyIGNsYXNzbmFtZSA9IFwiXCI7XG4gICAgICAgIGlmIChydW4uVG90YWwgPT0gcnVuLkNvbXBsZXRlKSB7XG4gICAgICAgICAgY2xhc3NuYW1lID0gXCJjb21wbGV0ZVwiO1xuICAgICAgICB9XG4gICAgICAgIGNlbGxzLnB1c2goPHRkIGtleT17XCJydW4tXCIgKyByaWR9IGNsYXNzTmFtZT17Y2xhc3NuYW1lfT57cnVuLkNvbXBsZXRlfSAvIHtydW4uVG90YWx9PC90ZD4pO1xuICAgICAgfVxuICAgICAgcm93cy5wdXNoKDx0ciBrZXk9e1wid2YtXCIgKyB3ZmlkfT57Y2VsbHN9PC90cj4pO1xuICAgIH1cblxuICAgIFJlYWN0RE9NLnJlbmRlcihcbiAgICAgICg8ZGl2PlxuICAgICAgICA8aDM+TW9ydGFyPC9oMz5cbiAgICAgICAgPHRhYmxlPlxuICAgICAgICAgIDx0aGVhZD48dHI+e2hlYWRlcn08L3RyPjwvdGhlYWQ+XG4gICAgICAgICAgPHRib2R5Pntyb3dzfTwvdGJvZHk+XG4gICAgICAgIDwvdGFibGU+XG4gICAgICA8L2Rpdj4pLFxuICAgICAgZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ3Jvb3QnKVxuICAgICk7XG59KX0pO1xuXG5SZWFjdERPTS5yZW5kZXIoXG4gICg8ZGl2PlxuICAgIDxoMz5Nb3J0YXI8L2gzPjxkaXY+TG9hZGluZzwvZGl2PlxuICA8L2Rpdj4pLFxuICBkb2N1bWVudC5nZXRFbGVtZW50QnlJZCgncm9vdCcpXG4pO1xuIiwiZnVuY3Rpb24gRm9vKCkge1xuICByZXR1cm4gXCJGb29cIlxufVxuXG5leHBvcnQgeyBGb28gfTtcbiJdfQ==
