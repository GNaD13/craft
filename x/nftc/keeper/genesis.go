package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/craft/x/nftc/types"
)

// InitGenesis new nft genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	for _, class := range data.Classes {
		if err := k.SaveClass(ctx, *class); err != nil {
			panic(err)
		}
	}
	for _, entry := range data.Entries {
		for _, nft := range entry.Nfts {
			owner, err := sdk.AccAddressFromBech32(entry.Owner)
			if err != nil {
				panic(err)
			}

			if err := k.Mint(ctx, *nft, owner); err != nil {
				panic(err)
			}
		}
	}
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	classes := k.GetClasses(ctx)
	nftMap := make(map[string][]*types.NFT)
	for _, class := range classes {
		nfts := k.GetNFTsOfClass(ctx, class.Id)
		for i, n := range nfts {
			owner := k.GetOwner(ctx, n.ClassId, n.Id)
			nftArr, ok := nftMap[owner.String()]
			if !ok {
				nftArr = make([]*types.NFT, 0)
			}
			nftMap[owner.String()] = append(nftArr, &nfts[i])
		}
	}

	owners := make([]string, 0, len(nftMap))
	for owner := range nftMap {
		owners = append(owners, owner)
	}
	sort.Strings(owners)

	entries := make([]*types.Entry, 0, len(nftMap))
	for _, owner := range owners {
		entries = append(entries, &types.Entry{
			Owner: owner,
			Nfts:  nftMap[owner],
		})
	}
	return &types.GenesisState{
		Classes: classes,
		Entries: entries,
	}
}
