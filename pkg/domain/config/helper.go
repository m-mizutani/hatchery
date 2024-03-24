package config

func (x Config) LookupAction(id string) Action {
	for _, action := range x.Actions {
		if action.GetId() == id {
			return action
		}
	}

	return nil
}
