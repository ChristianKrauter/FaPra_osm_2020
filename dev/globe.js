$.ajax({
    url: "/grid",
    data: data
}).done(function(data) {
    if (data == "false") {

    } else {
        points = JSON.parse(data)
        for (i = 0; i < points.length; i++) {
            for (j = 0; j < points[i].length; j++) {
                point = points[i][j]
                createRedPoint(Cesium.Cartesian3.fromDegrees(point[1], point[0]))
                //createPoint(point[0],point[1])
            }
        }
    }

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
/*createRedPoint(Cesium.Cartesian3.fromDegrees(0.0, 0.0))
$.ajax({
    url: "/testneighbours",
    data: { lat: "0.0", lng: "0.0" }
}).done(function(response) {
    if (data == "false") {

    } else {
        console.log(response)
        var points = JSON.parse(response)
        //
        for (i = 0; i < points.length; i++) {
            point = points[i]
            createPoint(Cesium.Cartesian3.fromDegrees(point[1], point[0]));
        }
    }

});*/

if (!viewer.scene.pickPositionSupported) {
    window.alert("This browser does not support pickPosition.");
}

viewer.cesiumWidget.screenSpaceEventHandler.removeInputAction(
    Cesium.ScreenSpaceEventType.LEFT_DOUBLE_CLICK
);

function createPoint(worldPosition) {
    var point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: Cesium.Color.WHITE,
            pixelSize: 8,
            heightReference: Cesium.HeightReference.CLAMP_TO_GROUND,
        },
    });
    return point;
}

function createRedPoint(worldPosition) {
    var point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: Cesium.Color.RED,
            pixelSize: 4,
            heightReference: Cesium.HeightReference.CLAMP_TO_GROUND,
        },
    });
    return point;
}

function drawLine(viewer, coords) {
    viewer.entities.add({
        polyline: {
            positions: Cesium.Cartesian3.fromDegreesArray(coords),
            width: 1,
            material: Cesium.Color.DEEPSKYBLUE,
        },
        label: {
            position: coords.slice(0, 2),
            text: 'Select points on water',
            font: '14pt monospace',
            style: Cesium.LabelStyle.FILL_AND_OUTLINE,
            outlineWidth: 2,
            verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
        },
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

function onLeftMouseClick(event) {
    // We use `viewer.scene.pickPosition` here instead of `viewer.camera.pickEllipsoid` so that
    // we get the correct point when mousing over terrain.
    var earthPosition = viewer.scene.pickPosition(event.position);

    // `earthPosition` will be undefined if our mouse is not over the globe.
    if (Cesium.defined(earthPosition)) {

        const cartographic = viewer.scene.globe.ellipsoid.cartesianToCartographic(earthPosition);
        const longitudeString = Cesium.Math.toDegrees(cartographic.longitude).toFixed(15);
        const latitudeString = Cesium.Math.toDegrees(cartographic.latitude).toFixed(15);
        //console.log("1: ")
        //console.log(earthPosition)

        /*var data = {
            lat: latitudeString,
            lng: longitudeString,
        }
        $.ajax({
            url: "/testpoint",
            data: data
        }).done(function(response) {
            if (data == "false") {

            } else {
                console.log(response)
                var point = JSON.parse(response)
                createPoint(Cesium.Cartesian3.fromDegrees(point[0], point[1]));
            }

        });*/

        if (data["startLat"] == "") {
            data["startLat"] = latitudeString
            data["startLng"] = longitudeString
            createPoint(earthPosition);
        } else {

            data["endLat"] = latitudeString
            data["endLng"] = longitudeString

            $.ajax({
                url: "/point",
                data: data
            }).done(function(data) {
                if (data == "false") {

                } else {
                    createPoint(earthPosition);
                    var features = JSON.parse(data).features
                    var coord1d = []
                    for (i = 0; i < features.length; i++) {
                        var coordinates = features[i].geometry.coordinates;
                        coordinates.forEach(element => coord1d.push(element[1], element[0]))
                    }
                    drawLine(viewer, coord1d)
                }

            });

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
                window.alert(
                    "This browser does not support polylines on terrain."
                );
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