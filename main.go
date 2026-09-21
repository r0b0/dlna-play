package main

import (
	"encoding/xml"
	"fmt"
	"github.com/huin/goupnp/dcps/av1"
	// "strings"
)

type DIDLLite struct {
	XMLName    xml.Name    `xml:"DIDL-Lite"`
	Containers []Container `xml:"container"`
	Items      []Item      `xml:"item"`
}

type Container struct {
	ID       string `xml:"id,attr"`
	ParentID string `xml:"parentID,attr"`
	Title    string `xml:"title"`
}

type Item struct {
	ID        string     `xml:"id,attr"`
	ParentID  string     `xml:"parentID,attr"`
	Title     string     `xml:"title"`
	Class     string     `xml:"class"`
	Resources []Resource `xml:"res"`
}

type Resource struct {
	ProtocolInfo string `xml:"protocolInfo,attr"`
	URL          string `xml:",chardata"`
}

type UserItem struct {
	Directory *av1.ContentDirectory1
	Item      Item
	Name      string
}

func main() {
	/*
		devices, err := goupnp.DiscoverDevices("upnp:rootdevice")
		if err != nil {
			fmt.Println(err)
		}
		if len(devices) > 0 {
			for _, d := range devices {
				fmt.Printf("Discovered device %s\n", d.Root.Device.String())
				ListServices(d.Root.Device.Services)
				ListDevices(d.Root.Device.Devices)
			}
		}
	*/

	clients, _, err := av1.NewAVTransport1Clients()
	if err != nil {
		panic(err)
	}
	for _, c := range clients {
		fmt.Println("Found renderer:", c.Location)
	}

	directories, _, err := av1.NewContentDirectory1Clients()
	if err != nil {
		panic(err)
	}
	for _, d := range directories {
		fmt.Println("Found content directory:", d.RootDevice.Device.FriendlyName)
		// result, nr, matches, update_id, err := d.Browse("0", "BrowseDirectChildren", "*", 0, 100, "")
		items, err := BrowseRecursive(d, "0", "")
		if err != nil {
			panic(err)
		}
        for _, item := range(items) {
            fmt.Printf("  %s (%s)\n", item.Name, item.Item.Resources[0].URL)
        }
		// fmt.Printf("Browse Result: %s NumberReturned: %d TotalMatches: %d, UpdateID: %d\n", result, nr, matches, update_id)
	}
}

func BrowseRecursive(d *av1.ContentDirectory1, objectId string, prefix string) ([]UserItem, error) {
	itemsSoFar := []UserItem{}

	resultXML, _, _, _, err := d.Browse(objectId, "BrowseDirectChildren", "*", 0, 100, "")
	if err != nil {
		return itemsSoFar, err
	}
	var didl DIDLLite

	err = xml.Unmarshal([]byte(resultXML), &didl)
	if err != nil {
		return itemsSoFar, err
	}

	for _, c := range didl.Containers {
		// fmt.Printf("%sContainer: %s (%s)\n", prefix, c.Title, c.ID)
		switch c.Title {
		case "Music", "Pictures", "Video":
			// fmt.Printf("%s Ignoring\n", prefix)
		default:
			newItems, err := BrowseRecursive(d, c.ID, fmt.Sprintf("%s/%s", prefix, c.Title))
			if err != nil {
				return itemsSoFar, err
			}
            itemsSoFar = append(itemsSoFar, newItems...)
		}

	}

	for _, item := range didl.Items {
		// fmt.Printf("%sItem: %s (%s)\n", prefix, item.Title, item.ID)
		itemsSoFar = append(itemsSoFar, UserItem{
			Directory: d,
			Item:      item,
			Name:      fmt.Sprintf("%s/%s", prefix, item.Title),
		})

        /*
		for _, res := range item.Resources {
			fmt.Printf("%s URL: %s\n", prefix, strings.TrimSpace(res.URL))
			// fmt.Printf("%s Protocol: %s\n", prefix, res.ProtocolInfo)
		}
        */
	}

	return itemsSoFar, nil
}

/*

func ListDevices(devices []goupnp.Device) {
	if len(devices) > 0 {
		for _, sd := range devices {
			fmt.Printf("Discovered Device %s\n", sd.String())
			ListDevices(sd.Devices)
			ListServices(sd.Services)
		}
	}
}

func ListServices(services []goupnp.Service) {
	if len(services) > 0 {
		for _, s := range services {
			scpd, err := s.RequestSCDP()
			fmt.Printf("  Discovered service %s\n", s.ServiceId)
			if err != nil {
				fmt.Printf("    Error requesting service SCPD: %v\n", err)
			} else {
				fmt.Println("    Available actions:")
				for _, action := range scpd.Actions {
					fmt.Printf("    * %s\n", action.Name)
					for _, arg := range action.Arguments {
						var varDesc string
						if stateVar := scpd.GetStateVariable(arg.RelatedStateVariable); stateVar != nil {
							varDesc = fmt.Sprintf(" (%s)", stateVar.DataType.Name)
						}
						fmt.Printf("      * [%s] %s%s\n", arg.Direction, arg.Name, varDesc)
					}
				}
			}
		}
	}

}
*/
