//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var testData = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":160,"Age":38,"Name":"Jill Lloyd","Email":"jilllloyd@ozon.ru","Phone":"+1 (873) 450-3650","Address":"880 Abbey Court, Rodman, Michigan, 6324"}
{"Id":170,"Age":37,"Name":"Olsen Guthrie","Email":"olsenguthrie@mail.ru","Phone":"+1 (805) 458-3986","Address":"745 Congress Street, Silkworth, South Carolina, 7356"}
{"Id":197,"Age":24,"Name":"Mari Ballard","Email":"mariballard@yandex.RU","Phone":"+1 (830) 507-2834","Address":"228 Walker Court, Irwin, Connecticut, 6631"}
{"Id":169,"Age":21,"Name":"Sosa Norris","Email":"sosanorris@yandex.Ru","Phone":"+1 (918) 567-3757","Address":"566 Lincoln Place, Fresno, Texas, 3823"}
{"Id":113,"Age":30,"Name":"Aida Salinas","Email":"aidasalinas.us","Phone":"+1 (807) 484-2758","Address":"682 Tilden Avenue, Tyhee, Nevada, 1449"}
{"Id":111,"Age":50,"Name":"Joe Doe"}
{"Id":112,"Age":18,"Name":"Jane Doe", "Email": "jane@doe.mail.org"}
<notjson>Someone forgot an XML here =(</notjson>`

func TestGetDomainStat(t *testing.T) {
	t.Run("find 'mail.org' domain with 2 layers", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "mail.org")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"doe.mail.org": 1,
		}, result)
	})

	t.Run("find 'us' with incorrect email", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "us")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("find 'ru' with capital case in domain", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "ru")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"mail.ru":   1,
			"ozon.ru":   1,
			"yandex.ru": 2,
		}, result)
	})

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(testData), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}
