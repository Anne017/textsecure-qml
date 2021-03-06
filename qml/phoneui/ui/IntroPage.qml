import QtQuick 2.0
import Ubuntu.Components 0.1
import "../components"

TelegramPage {
    id: page

    head.backAction.visible: false

    body: Item {
        anchors {
            fill: parent
            margins: units.gu(2)
        }

        Text {
            id: infoText
            elide: Text.ElideRight
            anchors {
                top: parent.top
                margins: units.gu(1)
            }
            width: parent.width
            wrapMode: Text.WordWrap
            text: "<h3>Thanks for trying out TextSecure!</h3><br><br> \
            File bugs and feature requests on github:<br>\
            <a href='https://github.com/janimo/textsecure-qml/issues'>https://github.com/janimo/textsecure-qml/issues</a><br>"
            onLinkActivated:Qt.openUrlExternally(link)
        }

        TelegramButton {
            anchors {
                top: infoText.bottom
                topMargin: units.gu(1)
                right: parent.right
                left: parent.left
            }
            width: parent.width

            text: i18n.tr("OK")
            onClicked: pageStack.push(dialogsPage)
        }
    }
}
