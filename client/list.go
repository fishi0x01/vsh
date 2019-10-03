package client

func (client *Client) listTopLevel() (result []string) {
	for k := range client.KVBackends {
		result = append(result, k)
	}
	return result
}

func (client *Client) listLowLevel(path string) (result []string, err error) {
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return result, err
	}

	if s != nil {
		if keysInterface, ok := s.Data["keys"]; ok {
			for _, valInterface := range keysInterface.([]interface{}) {
				val := valInterface.(string)
				result = append(result, val)
			}
		}
	}

	return result, err
}