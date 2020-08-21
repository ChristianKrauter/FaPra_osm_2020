$("#clearButton").click(function() {
    viewer.entities.removeAll();
});

$("#gridButton").click(function() {
    $.ajax({
        url: "/basicGrid",
        data: data
    }).done(function(data) {
        if (data != "false") {
            points = JSON.parse(data)
            console.log(points)
            for (i = 0; i < points.length; i++) {
                createColoredPoint(Cesium.Cartesian3.fromDegrees(points[i][0], points[i][1]), Cesium.Color.RED, 4)
            }
        }
    });
});

var data = {
    startLat: "",
    startLng: "",
    endLat: "",
    endLng: ""
}

var drawingMode = "line";

var osm = new Cesium.OpenStreetMapImageryProvider({
    url: 'https://a.tile.openstreetmap.org/'
});

var viewer = new Cesium.Viewer("cesiumContainer", {
    selectionIndicator: false,
    infoBox: false,
    terrainProvider: Cesium.createWorldTerrain(),
});

if (!viewer.scene.pickPositionSupported) {
    window.alert("This browser does not support pickPosition.");
}

viewer.cesiumWidget.screenSpaceEventHandler.removeInputAction(
    Cesium.ScreenSpaceEventType.LEFT_DOUBLE_CLICK
);

function createPoint(worldPosition, processed = false, start = false) {
    var point
    if (processed) {
        createColoredPoint(worldPosition, Cesium.Color.RED, 4)
    } else {
        var text = "End"
        console.log(worldPosition)
        point = viewer.entities.add({
            position: worldPosition,
            point: {
                color: Cesium.Color.WHITE,
                pixelSize: 5,
                heightReference: Cesium.HeightReference.CLAMP_TO_GROUND,
            },
        });

        if (start == true) {
            text = "Start"
        }

        point = viewer.entities.add({
            position: worldPosition,
            label: {
                height: 20000000,
                text: text,
                font: '14pt monospace',
                style: Cesium.LabelStyle.FILL_AND_OUTLINE,
                outlineWidth: 2,
                verticalOrigin: Cesium.VerticalOrigin.TOP,
                pixelOffset: new Cesium.Cartesian2(0, 32),
                eyeOffset: new Cesium.Cartesian3(0, 0, -3000000)
            }
        });
    }
    return point;
}

function drawLine(viewer, coords) {
    viewer.entities.add({
        positions: coords.slice(0, 2),
        polyline: {
            positions: Cesium.Cartesian3.fromDegreesArray(coords),
            width: 2,
            material: Cesium.Color.DEEPSKYBLUE,
            // granularity: Cesium.Math.toRadians(0.001)
            // distanceDisplayCondition: Cesium.Cartesian3(0, 0, 100000000)
        }
    })
}

function drawShape(positionData) {
    var shape;
    if (drawingMode === "line") {
        shape = viewer.entities.add({
            polyline: {
                positions: positionData,
                clampToGround: true,
                width: 3,
            },
        });
    } else if (drawingMode === "polygon") {
        shape = viewer.entities.add({
            polygon: {
                hierarchy: positionData,
                material: new Cesium.ColorMaterialProperty(
                    Cesium.Color.WHITE.withAlpha(0.7)
                ),
            },
        });
    }
    return shape;
}

function dijkstraProcessing(jsonData) {
    if (jsonData == "false") {
        //Todo
    } else {
        createPoint(this.earthPosition.ep);
        var features = JSON.parse(jsonData).features
        var coord1d = []
        for (i = 0; i < features.length; i++) {
            var coordinates = features[i].geometry.coordinates;
            coordinates.forEach(element => coord1d.push(element[0], element[1]))
        }
        drawLine(viewer, coord1d)
    }
}

function dijkstraAllNodesProcessing(jsonData) {
    if (jsonData == "false") {
        //Todo
    } else {
        tempData = JSON.parse(jsonData)
        route = tempData.Route
        processedNodes = tempData.AllNodes
        createPoint(this.earthPosition.ep);
        var features = route.features
        var coord1d = []
        for (i = 0; i < features.length; i++) {
            var coordinates = features[i].geometry.coordinates;
            coordinates.forEach(element => coord1d.push(element[0], element[1]))
        }
        drawLine(viewer, coord1d)
        for (i = 0; i < processedNodes.length; i++) {
            point = processedNodes[i]
            createPoint(Cesium.Cartesian3.fromDegrees(point[0], point[1]), true)
        }
    }
}

function onLeftMouseClick(event) {
    // We use `viewer.scene.pickPosition` here instead of `viewer.camera.pickEllipsoid` so that
    // we get the correct point when mousing over terrain.
    var earthPosition = viewer.scene.pickPosition(event.position);

    // `earthPosition` will be undefined if our mouse is not over the globe.
    if (Cesium.defined(earthPosition)) {

        const cartographic = viewer.scene.globe.ellipsoid.cartesianToCartographic(earthPosition);
        const longitudeString = Cesium.Math.toDegrees(cartographic.longitude).toFixed(15);
        const latitudeString = Cesium.Math.toDegrees(cartographic.latitude).toFixed(15);

        if (data["startLat"] == "") {
            data["startLat"] = latitudeString
            data["startLng"] = longitudeString
            createPoint(earthPosition, false, true);
        } else {
            data["endLat"] = latitudeString
            data["endLng"] = longitudeString

            if (document.getElementById("processedNodes").checked) {
                $.ajax({
                    url: "/dijkstraAllNodes",
                    data: data,
                    earthPosition: { ep: earthPosition }
                }).done(dijkstraAllNodesProcessing);
            } else {
                $.ajax({
                    url: "/dijkstra",
                    data: data,
                    earthPosition: { ep: earthPosition }
                }).done(dijkstraProcessing);
            }
            data = {
                startLat: "",
                startLng: "",
                endLat: "",
                endLng: ""
            }
        }
    }
}

var activeShapePoints = [];
var activeShape;
var floatingPoint;
var handler = new Cesium.ScreenSpaceEventHandler(viewer.canvas);
handler.setInputAction(onLeftMouseClick, Cesium.ScreenSpaceEventType.LEFT_CLICK);

// Redraw the shape so it's not dynamic and remove the dynamic shape.
function terminateShape() {
    activeShapePoints.pop();
    drawShape(activeShapePoints);
    viewer.entities.remove(floatingPoint);
    viewer.entities.remove(activeShape);
    floatingPoint = undefined;
    activeShape = undefined;
    activeShapePoints = [];
}

var options = [{
        text: "Draw Lines",
        onselect: function() {
            if (!Cesium.Entity.supportsPolylinesOnTerrain(viewer.scene)) {
                window.alert("This browser does not support polylines on terrain.");
            }
            terminateShape();
            drawingMode = "line";
        },
    },
    {
        text: "Draw Polygons",
        onselect: function() {
            terminateShape();
            drawingMode = "polygon";
        },
    },
];

function createColoredPoint(worldPosition, color, size) {
    var point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: color,
            pixelSize: size,
            heightReference: Cesium.HeightReference.CLAMP_TO_GROUND,
        },
    });
    return point;
}

// Zoom in to an area with mountains
viewer.camera.lookAt(
    //Cesium.Cartesian3.fromDegrees(-122.2058, 46.1955, 1000.0),
    //new Cesium.Cartesian3(5000.0, 5000.0, 5000.0)
    Cesium.Cartesian3.fromDegrees(0, 0, 0.0),
    new Cesium.Cartesian3(0.0, 0.0, 42000000.0)
);

viewer.camera.lookAtTransform(Cesium.Matrix4.IDENTITY);
viewer.timeline.destroy()
viewer.sceneModePicker.destroy()
viewer.navigationHelpButton.destroy()
//viewer.baseLayerPicker.destroy()
viewer.homeButton.destroy()
viewer.geocoder.destroy()
viewer.animation.destroy()