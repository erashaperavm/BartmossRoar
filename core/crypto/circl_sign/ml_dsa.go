package circl_sign

import (
	"github.com/cloudflare/circl/sign/mldsa/mldsa44"
	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
)

/*
	use CIRCL（cloudflare/circl/sign/mldsa）
	ml-dsa-44
	ml-dsa-65
	ml-dsa-87
*/

type mldsa44Pubkey mldsa44.PublicKey
type mldsa65PubKey mldsa65.PublicKey
type mldsa87PubKey mldsa87.PublicKey
