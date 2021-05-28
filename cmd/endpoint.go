package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/cfg"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type endpointCmd struct {
	cmd *cobra.Command
	cfg *cfg.Config
	sc  *svix.Svix
}

func newEndpointCmd(cfg *cfg.Config, s *svix.Svix) *endpointCmd {
	ec := &endpointCmd{
		sc:  s,
		cfg: cfg,
	}
	ec.cmd = &cobra.Command{
		Use:   "endpoint",
		Short: "List, create & modify endpoints",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {

			l, err := s.Application.List(getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.PrintApplicationList(l)
			return nil
		},
	}
	addFilterFlags(list)
	ec.cmd.AddCommand(list)

	// create

	create := &cobra.Command{
		Use:   "create APP_ID URL VERSION [DESCRIPTION] [FILTER_TYPE FILTER_TYPE ...]",
		Short: "Create a new endpoint",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			url := args[1]

			version, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return fmt.Errorf("version must be a valid int32")
			}

			var desc *string
			if len(args) >= 3 {
				desc = &args[1]
			}

			var filterTypes *[]string
			if len(args) >= 4 {
				filters := args[4:]
				filterTypes = &filters
			}

			ep := &svix.EndpointIn{
				Url:         url,
				Version:     int32(version),
				Description: desc,
				FilterTypes: filterTypes,
			}

			out, err := s.Endpoint.Create(appID, ep)
			if err != nil {
				return err
			}
			fmt.Println("Endpoint Created!")
			pretty.PrintEndpointOut(out)
			return nil
		},
	}
	ec.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID ENDPOINT_ID",
		Short: "Get an endpoint by id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			endpointID := args[1]

			out, err := s.Endpoint.Get(appID, endpointID)
			if err != nil {
				return err
			}

			pretty.PrintEndpointOut(out)
			return nil
		},
	}
	ec.cmd.AddCommand(get)

	update := &cobra.Command{
		Use:   "update APP_ID NAME [UID]",
		Short: "Update an application by id",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			name := args[1]
			var uid *string
			if len(args) >= 2 {
				uid = &args[2]
			}

			app := &svix.ApplicationIn{
				Name: name,
				Uid:  uid,
			}

			out, err := s.Application.Update(appID, app)
			if err != nil {
				return err
			}

			pretty.PrintApplicationOut(out)
			return nil
		},
	}
	ec.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete APP_ID ENDPOINT_ID",
		Short: "Delete an endpoint by id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			err := s.Endpoint.Delete(appID, endpointID)
			if err != nil {
				return err
			}

			fmt.Printf("Endpoint \"%s\" Deleted!\n", endpointID)
			return nil
		},
	}
	ec.cmd.AddCommand(delete)

	secret := &cobra.Command{
		Use:   "secret APP_ID ENDPOINT_ID",
		Short: "get an endpoint's secret by id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			out, err := s.Endpoint.GetSecret(appID, endpointID)
			if err != nil {
				return err
			}

			pretty.PrintEndpointSecret(endpointID, out)
			return nil
		},
	}
	ec.cmd.AddCommand(secret)

	return ec
}

func (ec *endpointCmd) Cmd() *cobra.Command {
	return ec.cmd
}
