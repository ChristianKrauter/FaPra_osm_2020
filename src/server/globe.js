var viewer = new Cesium.Viewer("cesiumContainer", {
    selectionIndicator: false,
    infoBox: false,
    terrainProvider: Cesium.createWorldTerrain(),
});
// Zoom in to an area with mountains
viewer.camera.lookAt(
    Cesium.Cartesian3.fromDegrees(0, 0, 0.0),
    new Cesium.Cartesian3(0.0, 0.0, 42000000.0)
);

var handler = new Cesium.ScreenSpaceEventHandler(viewer.canvas);
handler.setInputAction(onLeftMouseClick, Cesium.ScreenSpaceEventType.LEFT_CLICK);
viewer.camera.lookAtTransform(Cesium.Matrix4.IDENTITY);
viewer.timeline.destroy()
viewer.sceneModePicker.destroy()
viewer.navigationHelpButton.destroy()
viewer.homeButton.destroy()
viewer.geocoder.destroy()
viewer.animation.destroy()

var routeData = {
    startLat: "",
    startLng: "",
    endLat: "",
    endLng: ""
}

if (!viewer.scene.pickPositionSupported) {
    window.alert("This browser does not support pickPosition.");
}

viewer.cesiumWidget.screenSpaceEventHandler.removeInputAction(
    Cesium.ScreenSpaceEventType.LEFT_DOUBLE_CLICK
);

function createPoint(worldPosition, processed = false, start = false) {
    var point
    if (processed) {
        var x = document.getElementById("algos").value
        switch (x) {
            case "0":
                createColoredPoint(worldPosition, Cesium.Color.DARKSALMON, 2)
                break
            case "1":
                createColoredPoint(worldPosition, Cesium.Color.LIGHTGREEN , 2)
                break
            case "2":
                createColoredPoint(worldPosition, Cesium.Color.KHAKI, 2)
                break
            case "3":
                createColoredPoint(worldPosition, Cesium.Color.LAVENDER, 2)
                break
            case "4":
                createColoredPoint(worldPosition, Cesium.Color.LIGHTCYAN, 4)
                break
        }
    } else {
        var text = "End"
        point = viewer.entities.add({
            position: worldPosition,
            point: {
                color: Cesium.Color.WHITE,
                pixelSize: 5,
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
    var color = Cesium.RED
    var x = document.getElementById("algos").value
    var width = 4
    switch (x) {
        case "0":
            color = Cesium.Color.CRIMSON
            break
        case "1":
            color = Cesium.Color.GREENYELLOW
            break
        case "2":
            color = Cesium.Color.ORANGE
            break
        case "3":
            color = Cesium.Color.PURPLE
            break
        case "4":
            color = Cesium.Color.ROYALBLUE
            width = 2
            break
    }
    viewer.entities.add({
        positions: coords.slice(0, 2),
        polyline: {
            positions: Cesium.Cartesian3.fromDegreesArray(coords),
            width: width,
            material: color,

        }
    })
}

function createColoredPoint(worldPosition, color, size) {
    var point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: color,
            pixelSize: size,
        },
    });
    return point;
}

$("#clearButton").click(function() {
    viewer.entities.removeAll();
});

$("#gridButton").click(function() {
    $.ajax({
        url: "/Grid",
    }).done(function(data) {
        if (data != "false") {
            points = JSON.parse(data)
            for (i = 0; i < points.length; i++) {
                createColoredPoint(Cesium.Cartesian3.fromDegrees(points[i][0], points[i][1]), Cesium.Color.RED, 4)
            }
        }
    });
});

function routeProcessing(jsonData) {
    if (jsonData == "false") {
        //Todo
    } else {
        tempData = JSON.parse(jsonData).features
        console.log(tempData)
        for (i = 0; i < tempData.length; i++) {
            var coord1d = []
            var coordinates = tempData[i].geometry.coordinates;
            coordinates.forEach(element => coord1d.push(element[0], element[1]))
            drawLine(viewer, coord1d)
        }
    }
}

function routeAndNodesProcessing(jsonData) {
    if (jsonData == "false") {
        //Todo
    } else {
        tempData = JSON.parse(jsonData)
        console.log(tempData)
        for (i = 0; i < tempData.Route.features.length; i++) {
            var coord1d = []
            var coordinates = tempData.Route.features[i].geometry.coordinates;
            coordinates.forEach(element => coord1d.push(element[0], element[1]))
            drawLine(viewer, coord1d)
        }
        for (i = 0; i < tempData.AllNodes.length; i++) {
            createPoint(Cesium.Cartesian3.fromDegrees(tempData.AllNodes[i][0], tempData.AllNodes[i][1]), true)
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
        const longitudeString = Cesium.Math.toDegrees(cartographic.longitude).toFixed(20);
        const latitudeString = Cesium.Math.toDegrees(cartographic.latitude).toFixed(20);
        routeData["algo"] = document.getElementById("algos").value
        if (routeData["startLat"] == "") {
            routeData["startLat"] = latitudeString
            routeData["startLng"] = longitudeString
            $.ajax({
                url: "/getGridPoint",
                data: {
                    "lng": longitudeString,
                    "lat": latitudeString
                },
                earthPosition: { ep: earthPosition },
                'success': function(data) {
                    if (data == "false") {
                        alert("Please select a point on water connected with the oceans.")
                        routeData = {
                            startLat: "",
                            startLng: "",
                            endLat: "",
                            endLng: ""
                        }
                    } else {
                        data = JSON.parse(data).coordinates
                        createPoint(Cesium.Cartesian3.fromDegrees(data[0], data[1]), false, true);
                    }
                }
            })
        } else {
            routeData["endLat"] = latitudeString
            routeData["endLng"] = longitudeString
            $.ajax({
                url: "/getGridPoint",
                data: {
                    "lng": longitudeString,
                    "lat": latitudeString
                },
                earthPosition: { ep: earthPosition },
                'success': function(data) {
                    if (data == "false") {
                        alert("Please select a point on water connected with the oceans.")
                        routeData["endLat"] = ""
                        routeData["endLng"] = ""
                    } else {
                        data = JSON.parse(data).coordinates
                        createPoint(Cesium.Cartesian3.fromDegrees(data[0], data[1]), false);
                        if (document.getElementById("processedNodes").checked) {
                            $.ajax({
                                url: "/wayfindingAllNodes",
                                data: routeData,
                                earthPosition: { ep: earthPosition }
                            }).done(routeAndNodesProcessing);
                        } else {
                            $.ajax({
                                url: "/wayfinding",
                                data: routeData,
                                earthPosition: { ep: earthPosition }
                            }).done(routeProcessing);
                        }
                        routeData = {
                            startLat: "",
                            startLng: "",
                            endLat: "",
                            endLng: ""
                        }
                    }
                }
            })
        }
    }
}